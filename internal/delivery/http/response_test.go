package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponder_JSON(t *testing.T) {
	rs, _ := newTestResponder()
	rec := httptest.NewRecorder()

	rs.JSON(rec, get("/"), http.StatusCreated, map[string]string{"name": "budi"})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", got)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"name":"budi"}` {
		t.Errorf("body = %q", got)
	}
}

func TestResponder_JSON_ValueThatCannotBeEncoded(t *testing.T) {
	rs, logs := newTestResponder()
	rec := httptest.NewRecorder()

	rs.JSON(rec, get("/"), http.StatusOK, make(chan int))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"code":"internal_error"`) {
		t.Errorf("body is not the error shape: %s", rec.Body.String())
	}

	if !strings.Contains(logs.String(), "encode response") {
		t.Errorf("the failure was not logged: %s", logs.String())
	}
}
