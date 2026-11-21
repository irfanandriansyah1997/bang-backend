package config

import (
	"errors"
	"log/slog"
	"os"
)

// Config holds the app settings read from environment variables.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	LogFile     string
	LogLevel    string
	AppEnv      string
	AppVersion  string
}

// Load reads the config from the real environment.
func Load() (Config, error) {
	return LoadFrom(os.Getenv)
}

// LoadFrom reads the config using getenv, so tests can pass a fake one.
func LoadFrom(getenv func(string) string) (Config, error) {
	cfg := Config{
		Port:        valueOr(getenv("PORT"), "8080"),
		DatabaseURL: valueOr(getenv("DATABASE_URL"), "postgres://app:app@localhost:5432/fieldservice?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET"),
		LogFile:     valueOr(getenv("LOG_FILE"), "logs/app.log"),
		LogLevel:    valueOr(getenv("LOG_LEVEL"), "info"),
		AppEnv:      valueOr(getenv("APP_ENV"), "local"),
		AppVersion:  valueOr(getenv("APP_VERSION"), "dev"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, errors.New("config: JWT_SECRET is required")
	}

	// The default URL is the local docker-compose login (app:app). Outside local,
	// a forgotten DATABASE_URL must fail at startup, not connect with a known password.
	if cfg.AppEnv != "local" && getenv("DATABASE_URL") == "" {
		return Config{}, errors.New("config: DATABASE_URL is required when APP_ENV is not local")
	}

	return cfg, nil
}

// LogValue controls how Config looks in logs. Secrets are left out on purpose.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("port", c.Port),
		slog.String("log_file", c.LogFile),
		slog.String("log_level", c.LogLevel),
	)
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
