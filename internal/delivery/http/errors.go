package http

import (
	"errors"
	"net/http"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// errorBody is the one error shape of the API: { "error": { "code": "...", "message": "..." } }.
type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (rs *Responder) Error(w http.ResponseWriter, r *http.Request, err error) {
	status, code, message := mapError(err)

	if status >= http.StatusInternalServerError {
		rs.log.Error(r.Context(), "request failed",
			logger.Payload{"error": err.Error()},
			logger.SDC{"error_code": code},
		)
	}

	rs.JSON(w, r, status, errorBody{Error: errorDetail{Code: code, Message: message}})
}

func mapError(err error) (status int, code, message string) {
	var validation *domain.ValidationError

	switch {
	case errors.As(err, &validation):
		return http.StatusUnprocessableEntity, "validation_failed", validation.Message

	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "not_found", "the resource was not found"

	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, "forbidden", "you are not allowed to do this"

	case errors.Is(err, domain.ErrInvalidTransition):
		return http.StatusConflict, "INVALID_STATUS_TRANSITION", "the status cannot change this way"

	default:
		return http.StatusInternalServerError, "internal_error", "something went wrong"
	}
}
