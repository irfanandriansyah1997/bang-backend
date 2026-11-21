package config

import (
	"errors"
	"log/slog"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	LogFile     string
	LogLevel    string
	AppEnv      string
	AppVersion  string
}

func Load() (Config, error) {
	return LoadFrom(os.Getenv)
}

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

	return cfg, nil
}

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
