package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
)

func TestResponder_Error(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{name: "not found", err: domain.ErrNotFound, wantStatus: 404, wantCode: "not_found", wantMessage: "the resource was not found"},
		{name: "forbidden", err: domain.ErrForbidden, wantStatus: 403, wantCode: "forbidden", wantMessage: "you are not allowed to do this"},
		{name: "invalid transition", err: domain.ErrInvalidTransition, wantStatus: 409, wantCode: "INVALID_STATUS_TRANSITION", wantMessage: "the status cannot change this way"},
		{name: "validation", err: &domain.ValidationError{Message: "at least one photo is required"}, wantStatus: 422, wantCode: "validation_failed", wantMessage: "at least one photo is required"},
		{name: "wrapped not found", err: fmt.Errorf("find work order 12: %w", domain.ErrNotFound), wantStatus: 404, wantCode: "not_found", wantMessage: "the resource was not found"},
		{name: "wrapped validation", err: fmt.Errorf("complete: %w", &domain.ValidationError{Message: "bad"}), wantStatus: 422, wantCode: "validation_failed", wantMessage: "bad"},
		{name: "unknown error", err: errors.New("boom"), wantStatus: 500, wantCode: "internal_error", wantMessage: "something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, _ := newTestResponder()
			rec := httptest.NewRecorder()

			rs.Error(rec, get("/"), tt.err)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var body errorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body is not the error shape: %v: %s", err, rec.Body.String())
			}

			if body.Error.Code != tt.wantCode || body.Error.Message != tt.wantMessage {
				t.Errorf("error = %+v, want code %q message %q", body.Error, tt.wantCode, tt.wantMessage)
			}
		})
	}
}

func TestResponder_Error_HidesDetailsOf5xx(t *testing.T) {
	rs, logs := newTestResponder()
	rec := httptest.NewRecorder()

	rs.Error(rec, get("/"), errors.New(`pq: password authentication failed for user "app"`))

	if strings.Contains(rec.Body.String(), "password") {
		t.Errorf("response leaks the error details: %s", rec.Body.String())
	}

	// The details are for us: logged once, with the error code in the indexed SDC.
	if got := strings.Count(logs.String(), "request failed"); got != 1 {
		t.Errorf("5xx was logged %d times, want 1: %s", got, logs.String())
	}

	if !strings.Contains(logs.String(), `"error_code":"internal_error"`) {
		t.Errorf("log line has no error_code in sdc: %s", logs.String())
	}
}

func TestResponder_Error_DoesNotLog4xx(t *testing.T) {
	rs, logs := newTestResponder()

	rs.Error(httptest.NewRecorder(), get("/"), domain.ErrNotFound)

	if logs.Len() != 0 {
		t.Errorf("a 404 must not be logged by the error mapper, got: %s", logs.String())
	}
}
