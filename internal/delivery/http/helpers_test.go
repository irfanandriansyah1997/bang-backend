package http

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// newTestLogger returns a logger that writes into the returned buffer, so a test can read the log lines.
func newTestLogger() (*logger.Logger, *bytes.Buffer) {
	var logs bytes.Buffer

	log := logger.New(&logs, logger.Options{Level: slog.LevelDebug, Service: "test", Env: "test", Version: "test"})

	return log, &logs
}

func newTestResponder() (*Responder, *bytes.Buffer) {
	log, logs := newTestLogger()

	return NewResponder(log), logs
}

func get(target string) *http.Request {
	return httptest.NewRequestWithContext(context.Background(), http.MethodGet, target, http.NoBody)
}
