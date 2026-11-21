package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"strings"
)

// Payload holds extra data for one log entry. It is stored under "payload"
// and is not indexed, so any value type is fine.
type Payload map[string]any

// SDC holds searchable fields for a log entry. It is stored under "sdc" and
// indexed in Kibana, so values are strings only.
type SDC map[string]string

// Extra is something that can be added to a log call: Payload or SDC.
type Extra interface {
	apply(e *entry)
}

type entry struct {
	payload Payload
	sdc     SDC
}

func (p Payload) apply(e *entry) { maps.Copy(e.payload, p) }
func (s SDC) apply(e *entry)     { maps.Copy(e.sdc, s) }

// Options configures a Logger.
type Options struct {
	Level   slog.Level
	Service string
	Env     string
	Version string
}

// Logger writes JSON log lines with service, env, and version on every line.
type Logger struct {
	base *slog.Logger
}

const timestampFormat = "2006-01-02T15:04:05.000Z07:00" // UTC, milliseconds

// New creates a Logger that writes JSON lines to w.
func New(w io.Writer, opts Options) *Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       opts.Level,
		ReplaceAttr: renameKeys,
	})

	base := slog.New(handler).With(
		slog.String("service", opts.Service),
		slog.String("env", opts.Env),
		slog.String("version", opts.Version),
	)

	return &Logger{base: base}
}

func renameKeys(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}

	switch a.Key {
	case slog.TimeKey:
		return slog.String("timestamp", a.Value.Time().UTC().Format(timestampFormat))
	case slog.MessageKey:
		return slog.String("message", a.Value.String())
	}

	return a
}

// Debug logs at debug level.
func (l *Logger) Debug(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelDebug, msg, extras)
}

// Info logs at info level.
func (l *Logger) Info(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelInfo, msg, extras)
}

// Warn logs at warn level.
func (l *Logger) Warn(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelWarn, msg, extras)
}

// Error logs at error level.
func (l *Logger) Error(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelError, msg, extras)
}

func (l *Logger) log(ctx context.Context, level slog.Level, msg string, extras []Extra) {
	if !l.base.Enabled(ctx, level) {
		return
	}

	e := entry{payload: Payload{}, sdc: SDC{}}
	maps.Copy(e.sdc, sdcFrom(ctx))
	for _, x := range extras {
		x.apply(&e)
	}

	var attrs []slog.Attr

	if id := RequestID(ctx); id != "" {
		attrs = append(attrs, slog.String("x_request_id", id))
	}

	if len(e.payload) > 0 {
		attrs = append(attrs, slog.Any("payload", map[string]any(e.payload)))
	}

	if len(e.sdc) > 0 {
		attrs = append(attrs, slog.Any("sdc", map[string]string(e.sdc)))
	}

	l.base.LogAttrs(ctx, level, msg, attrs...)
}

type sdcKey struct{}
type requestIDKey struct{}

// WithSDC returns a copy of ctx with sdc merged into the SDC already in it.
func WithSDC(ctx context.Context, sdc SDC) context.Context {
	merged := SDC{}
	maps.Copy(merged, sdcFrom(ctx))
	maps.Copy(merged, sdc)

	return context.WithValue(ctx, sdcKey{}, merged)
}

func sdcFrom(ctx context.Context) SDC {
	if sdc, ok := ctx.Value(sdcKey{}).(SDC); ok {
		return sdc
	}

	return nil
}

// WithRequestID returns a copy of ctx that carries the request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID returns the request ID stored in ctx, or "" if there is none.
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}

	return ""
}

// ParseLevel turns "debug", "info", "warn", or "error" into a slog.Level.
func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToUpper(s))); err != nil {
		return slog.LevelInfo, fmt.Errorf("logger: unknown level %q (use debug, info, warn, error)", s)
	}

	return level, nil
}

// OpenFile opens the log file for appending and creates its folders if needed.
func OpenFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, fmt.Errorf("logger: create log folder: %w", err)
	}

	f, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("logger: open log file: %w", err)
	}

	return f, nil
}
