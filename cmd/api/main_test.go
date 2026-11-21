package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/config"
)

const testSecret = "test-secret-do-not-log"

func testConfig(t *testing.T) config.Config {
	t.Helper()

	return config.Config{
		Port:        "8080",
		DatabaseURL: "postgres://app:app@127.0.0.1:1/fieldservice?sslmode=disable&connect_timeout=1",
		JWTSecret:   testSecret,
		LogFile:     filepath.Join(t.TempDir(), "logs", "app.log"),
		LogLevel:    "info",
		AppEnv:      "test",
		AppVersion:  "test",
	}
}

func TestRun_InvalidLogLevel(t *testing.T) {
	cfg := testConfig(t)
	cfg.LogLevel = "random level"

	err := run(context.Background(), cfg, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "log level") {
		t.Fatalf("run() error = %v, want an error about the log level", err)
	}
}

func TestRun_DatabaseUnreachable(t *testing.T) {
	err := run(context.Background(), testConfig(t), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "connect database") {
		t.Fatalf("run() error = %v, want an error about the database connection", err)
	}
}

func TestRun_ShutDownOnCancel(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		if os.Getenv("CI") == "true" {
			t.Fatalf("TEST_DATABASE_URL must be set in CI")
		}
		t.Skip("TEST_DATABASE_URL not set")
	}

	tests := []struct {
		name       string
		appEnv     string
		wantPretty bool
	}{
		{name: "local env indents the terminal", appEnv: "local", wantPretty: true},
		{name: "other env keeps one line per entry", appEnv: "production", wantPretty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testConfig(t)
			cfg.DatabaseURL = dbURL
			cfg.AppEnv = tt.appEnv

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var out bytes.Buffer
			done := make(chan error, 1)
			go func() { done <- run(ctx, cfg, &out) }()

			waitForLog(t, cfg.LogFile, "app started")
			cancel() // same as pressing Ctrl+C

			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("run() error: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("run() did not return after the context was canceled")
			}

			file, err := os.ReadFile(cfg.LogFile)
			if err != nil {
				t.Fatalf("read log file: %v", err)
			}

			lines := strings.Split(strings.TrimSpace(string(file)), "\n")
			if len(lines) != 3 {
				t.Errorf("log file has %d lines, want 3: %q", len(lines), file)
			}

			if !strings.Contains(string(file), `"message":"shutting down"`) {
				t.Errorf("log file is missing the shutdown message: %q", file)
			}

			if strings.Contains(string(file), testSecret) || strings.Contains(out.String(), testSecret) {
				t.Errorf("JWT secret was written to the logs")
			}

			if pretty := strings.Contains(out.String(), "\n  \"message\""); pretty != tt.wantPretty {
				t.Errorf("terminal output indented = %v, want %v:\n%s", pretty, tt.wantPretty, out.String())
			}
		})
	}
}

// waitForLog polls the log file until it contains want, or fails after 5 seconds.
func waitForLog(t *testing.T, path, want string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(path); err == nil && strings.Contains(string(data), want) {
			return
		}

		// Sleep is fine here: we poll a file written by another goroutine, and there is no signal to wait on.
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("log %q did not appear in %s within 5s", want, path)
}
