package utils

import (
	"errors"
	"fmt"
)

// Common error types
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrBadRequest    = errors.New("bad request")
	ErrInternal      = errors.New("internal server error")
	ErrAlreadyExists = errors.New("resource already exists")
)

// AppError represents an application error
type AppError struct {
	Err     error
	Message string
	Code    int
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized checks if an error is an unauthorized error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden checks if an error is a forbidden error
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsBadRequest checks if an error is a bad request error
func IsBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

// IsInternal checks if an error is an internal server error
func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}

// IsAlreadyExists checks if an error is an already exists error
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// WrapError wraps an error with a message
func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}