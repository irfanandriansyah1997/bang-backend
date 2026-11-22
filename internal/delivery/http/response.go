package http

import (
	"encoding/json"
	"net/http"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// Responder writes every JSON response of the API: normal bodies and errors.
// Handlers get one Responder and never write JSON themselves.
type Responder struct {
	log *logger.Logger
}

// NewResponder creates a Responder that logs through log.
func NewResponder(log *logger.Logger) *Responder {
	return &Responder{log: log}

}

// JSON writes v as a JSON body with the given status code.
func (rs *Responder) JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	// Encode first. Once WriteHeader is called, the status can no longer change,
	// so a value that cannot be encoded must be found before that.
	body, err := json.Marshal(v)
	if err != nil {
		rs.log.Error(r.Context(), "encode response", logger.Payload{"error": err.Error()})
		status = http.StatusInternalServerError
		body = []byte(`{"error":{"code":"internal_error","message":"something when wrong"}}`)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(append(body, '\n')); err != nil {
		rs.log.Warn(r.Context(), "write response", logger.Payload{"error": err.Error()})
	}
}
