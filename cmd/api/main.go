package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/config"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/database"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

const serviceName = "bang-backend"

func main() {
	os.Exit(realMain())
}

// realMain returns the exit code instead of calling os.Exit itself,
// so the deferred stop() below always runs.
func realMain() int {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// ctx is canceled when the user presses Ctrl+C (SIGINT) or the
	// platform sends SIGTERM (e.g. Docker stop).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

// run wires the app together and blocks until ctx is canceled.
func run(ctx context.Context, cfg config.Config, stdout io.Writer) error {
	level, err := logger.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parse log level: %w", err)
	}

	logFile, err := logger.OpenFile(cfg.LogFile)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	defer logFile.Close()

	// The terminal is indented only in local development.
	// The log file always gets the raw JSON line, one per entry.
	// The file comes first: MultiWriter stops at the first error, and the file is the full record.
	terminal := stdout
	if cfg.AppEnv == "local" {
		terminal = logger.NewPrettyWriter(stdout)
	}

	log := logger.New(io.MultiWriter(logFile, terminal), logger.Options{
		Level:   level,
		Service: serviceName,
		Env:     cfg.AppEnv,
		Version: cfg.AppVersion,
	})

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	log.Info(ctx, "database connected")
	log.Info(ctx, "app started", logger.Payload{"port": cfg.Port, "log_level": cfg.LogLevel})

	// The HTTP server will start here in #3. For now, wait for Ctrl+C.
	<-ctx.Done()

	log.Info(ctx, "shutting down")

	return nil
}
