package domain

import "errors"

var (
	// ErrNotFound means the thing that was asked for does not exist.
	ErrNotFound = errors.New("not found")

	// ErrForbidden means the user is known but is not allowed to do this.
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidTransition means a work order status cannot change that way.
	ErrInvalidTransition = errors.New("invalid status transition")
)

// ValidationError means the input broke a rule, for example
// "completing a work order needs at least one photo".
// Message is written for the person who sent the request.
type ValidationError struct {
	Message string
}

// Error makes *ValidationError an error.
func (e *ValidationError) Error() string {
	return "validation failed: " + e.Message
}
