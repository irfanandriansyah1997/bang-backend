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

type Payload map[string]any

type SDC map[string]string

type Extra interface {
	apply(e *entry)
}

type entry struct {
	payload Payload
	sdc     SDC
}

func (p Payload) apply(e *entry) { maps.Copy(e.payload, p) }
func (s SDC) apply(e *entry)     { maps.Copy(e.sdc, s) }

type Options struct {
	Level   slog.Level
	Service string
	Env     string
	Version string
}

type Logger struct {
	base *slog.Logger
}

const timestampFormat = "2006-01-02T15:04:05.000Z07:00" // UTC, milliseconds

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

func (l *Logger) Debug(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelDebug, msg, extras)
}

func (l *Logger) Info(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelInfo, msg, extras)
}

func (l *Logger) Warn(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelWarn, msg, extras)
}

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

func WithSDC(ctx context.Context, sdc SDC) context.Context {
	merged := SDC{}
	maps.Copy(merged, sdcFrom(ctx))
	maps.Copy(merged, sdc)

	return context.WithValue(ctx, sdcKey{}, merged)
}

func sdcFrom(ctx context.Context) SDC {
	sdc, _ := ctx.Value(sdcKey{}).(SDC)
	return sdc
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToUpper(s))); err != nil {
		return slog.LevelInfo, fmt.Errorf("logger: unknown level %q (use debug, info, warn, error)", s)
	}

	return level, nil
}

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
