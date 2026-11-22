package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestBusinessErrors_SurviceWrapping(t *testing.T) {
	tests := []struct {
		name   string
		target error
	}{
		{name: "not found", target: ErrNotFound},
		{name: "forbidden", target: ErrForbidden},
		{name: "invalid transition", target: ErrInvalidTransition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := fmt.Errorf("find work order 12: %w", tt.target)

			if !errors.Is(wrapped, tt.target) {
				t.Errorf("errors.is(%q, %q) = false, want true", wrapped, tt.target)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := fmt.Errorf("complete work order: %w", &ValidationError{Message: "at least one photo is required"})

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatal("errors.As did not find a *ValidationError")
	}

	if ve.Message != "at least one photo is required" {
		t.Errorf("Message = %q", ve.Message)
	}

	if got, want := ve.Error(), "validation failed: at least one photo is required"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
