package logger

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var testOpts = Options{Level: slog.LevelInfo, Service: "bang-backend", Env: "test", Version: "dev"}

func readLines(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var lines []map[string]any
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		var line map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			t.Fatalf("log line is not JSON %v\n%s", err, scanner.Text())
		}
		lines = append(lines, line)
	}

	return lines
}

func TestLog_Shapes(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, testOpts)
	ctx := context.Background()

	log.Info(ctx, "message")
	log.Info(ctx, "message", Payload{"user_id": 1, "sample_message": "lorem ipsum"})
	log.Info(ctx, "message", Payload{}, SDC{"correlation_id": "12345", "path": "/usr/asd/asd"})

	lines := readLines(t, &buf)
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}

	if _, ok := lines[0]["payload"]; ok {
		t.Errorf("line 1 must not have payload: %v", lines[0])
	}

	if _, ok := lines[0]["sdc"]; ok {
		t.Errorf("line 1 must not have sdc: %v", lines[0])
	}

	payload, _ := lines[1]["payload"].(map[string]any)
	if payload["user_id"] != float64(1) || payload["sample_message"] != "lorem ipsum" {
		t.Errorf("line 2 payload = %v", lines[1]["payload"])
	}

	if _, ok := lines[2]["payload"]; ok {
		t.Errorf("line 3 empty payload must be omitted: %v", lines[2])
	}
	sdc, _ := lines[2]["sdc"].(map[string]any)
	if sdc["correlation_id"] != "12345" || sdc["path"] != "/usr/asd/asd" {
		t.Errorf("line 3 sdc = %v", lines[2]["sdc"])
	}

	for i, line := range lines {
		if line["service"] != "bang-backend" || line["env"] != "test" || line["version"] != "dev" {
			t.Errorf("line %d missing service / env / version: %v", i+1, line)
		}
	}
}

func TestLog_KeynamesAndTimestamp(t *testing.T) {
	var buf bytes.Buffer
	New(&buf, testOpts).Info(context.Background(), "hello")

	line := readLines(t, &buf)[0]
	if line["message"] != "hello" {
		t.Errorf("message = %v, want hello", line["message"])
	}

	for _, old := range []string{"msg", "time"} {
		if _, ok := line[old]; ok {
			t.Errorf("old key %q must be renamed: %v", old, line)
		}
	}

	ts, _ := line["timestamp"].(string)
	if !strings.HasSuffix(ts, "Z") {
		t.Errorf("timestamp %q must be UTC (end with Z)", ts)
	}
	if _, err := time.Parse(timestampFormat, ts); err != nil {
		t.Errorf("timestamp %q has the wrong format: %v", ts, err)
	}
}

func TestLog_Levels(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, testOpts)
	ctx := context.Background()

	log.Debug(ctx, "hidden")
	log.Info(ctx, "info")
	log.Warn(ctx, "warn")
	log.Error(ctx, "error")

	lines := readLines(t, &buf)
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3 (debug logger must be dropped)", len(lines))
	}

	for i, want := range []string{"INFO", "WARN", "ERROR"} {
		if lines[i]["level"] != want {
			t.Errorf("line %d level = %v, want %s", i+1, lines[i]["level"], want)
		}
	}
}

func TestContext_RequestIDAndSDC(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, testOpts)

	ctx := WithRequestID(context.Background(), "b7f2")
	ctx = WithSDC(ctx, SDC{"user_id": "7", "role": "technician"})

	log.Info(ctx, "work order started", SDC{"work_order_id": "12", "role": "admin"})
	log.Info(context.Background(), "no context")

	lines := readLines(t, &buf)
	if lines[0]["x_request_id"] != "b7f2" {
		t.Errorf("x_request_id = %v, want b7f2", lines[0]["x_request_id"])
	}

	sdc, _ := lines[0]["sdc"].(map[string]any)
	if sdc["user_id"] != "7" || sdc["work_order_id"] != "12" || sdc["role"] != "admin" {
		t.Errorf("sdc = %v, want user_id=7, work_order_id=12, role=admin", sdc)
	}
	if _, ok := lines[1]["x_request_id"]; ok {
		t.Errorf("line 2 should be without context must not have x_request_id: %v", lines[1])
	}
}

func TestPrettyWriter(t *testing.T) {
	var buf bytes.Buffer
	New(NewPrettyWriter(&buf), testOpts).Info(context.Background(), "pretty")

	out := buf.String()
	if !strings.Contains(out, "\n  \"level\": \"INFO\"") {
		t.Errorf("output is not indented:\n%s", out)
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in      string
		want    slog.Level
		wantErr bool
	}{
		{in: "debug", want: slog.LevelDebug},
		{in: "info", want: slog.LevelInfo},
		{in: "WARN", want: slog.LevelWarn},
		{in: "error", want: slog.LevelError},
		{in: "verbose", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseLevel(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestOpenFile_CreatesFolder(t *testing.T) {
	// Two missing folders: os.Mkdir would fail here, os.MkdirAll works.
	path := filepath.Join(t.TempDir(), "nested", "deeper", "app.log")

	f, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile() error: %v", err)
	}
	defer f.Close()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("log file was not created: %v", err)
	}
}
