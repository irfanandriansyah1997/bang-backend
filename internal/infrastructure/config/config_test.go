package config

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func fakeEnv(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestLoadFrom(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		wantPort     string
		wantLogLevel string
		wantErr      bool
	}{
		{name: "defaults", env: map[string]string{"JWT_SECRET": "s3cret"}, wantPort: "8080", wantLogLevel: "info"},
		{name: "custom values", env: map[string]string{"JWT_SECRET": "s3cret", "PORT": "9090", "LOG_LEVEL": "debug"}, wantPort: "9090", wantLogLevel: "debug"},
		{name: "missing secret", env: map[string]string{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFrom(fakeEnv(tt.env))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Port != tt.wantPort {
				t.Errorf("Port = %q, want %q", cfg.Port, tt.wantPort)
			}

			if cfg.LogLevel != tt.wantLogLevel {
				t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, tt.wantLogLevel)
			}
		})
	}
}

func TestConfig_LogValue_HidesSecrets(t *testing.T) {
	cfg := Config{
		Port:        "8080",
		DatabaseURL: "postgres://app:supersecret@localhost:5432/db",
		JWTSecret:   "jwt-secret-value",
		LogFile:     "logs/app.log",
		LogLevel:    "info",
	}

	var buf bytes.Buffer
	slog.New(slog.NewJSONHandler(&buf, nil)).Info("config loaded", "config", cfg)
	out := buf.String()

	for _, secret := range []string{"supersecret", "jwt-secret-value"} {
		if strings.Contains(out, secret) {
			t.Errorf("log output leaks %q: %s", secret, out)
		}
	}

	if !strings.Contains(out, `"port":"8080"`) {
		t.Errorf("log output should include the port: %s", out)
	}
}
