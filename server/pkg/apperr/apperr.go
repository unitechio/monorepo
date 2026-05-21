// Package apperr defines typed application errors for the eEnglish platform.
//
// Usage pattern (in service layer):
//
//	return nil, apperr.NotFound("user", id)
//	return nil, apperr.Conflict("email already exists")
//	return nil, apperr.Internal(err)
//
// In the middleware, errors are unwrapped and serialized to a consistent JSON
// envelope with the HTTP status code, a human-readable message, and the
// request_id for traceability.
package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is a structured application error that carries an HTTP status code,
// a human-readable message, and an optional underlying cause.
type AppError struct {
	// Code is the HTTP status code that should be returned to the client.
	Code int `json:"-"`
	// Message is a safe, human-readable message for the API consumer.
	Message string `json:"message"`
	// Err is the internal cause, never exposed to the client.
	Err error `json:"-"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap enables errors.Is / errors.As chain traversal.
func (e *AppError) Unwrap() error {
	return e.Err
}

// ---------------------------------------------------------------------------
// Constructor helpers
// ---------------------------------------------------------------------------

// New creates a generic AppError.
func New(code int, message string, cause ...error) *AppError {
	var err error
	if len(cause) > 0 {
		err = cause[0]
	}
	return &AppError{Code: code, Message: message, Err: err}
}

// BadRequest (400) signals invalid client input.
func BadRequest(msg string, cause ...error) *AppError {
	return New(http.StatusBadRequest, msg, cause...)
}

// Unauthorized (401) signals missing or invalid credentials.
func Unauthorized(msg string, cause ...error) *AppError {
	return New(http.StatusUnauthorized, msg, cause...)
}

// Forbidden (403) signals the caller lacks permission.
func Forbidden(msg string, cause ...error) *AppError {
	return New(http.StatusForbidden, msg, cause...)
}

// NotFound (404) signals a missing resource.
func NotFound(resource string, id any) *AppError {
	return New(http.StatusNotFound, fmt.Sprintf("%s not found: %v", resource, id))
}

// Conflict (409) signals a uniqueness or state conflict (e.g., duplicate email).
func Conflict(msg string, cause ...error) *AppError {
	return New(http.StatusConflict, msg, cause...)
}

// UnprocessableEntity (422) signals semantically invalid data.
func UnprocessableEntity(msg string, cause ...error) *AppError {
	return New(http.StatusUnprocessableEntity, msg, cause...)
}

// TooManyRequests (429) signals rate-limit exceeded.
func TooManyRequests(msg string) *AppError {
	return New(http.StatusTooManyRequests, msg)
}

// Internal (500) wraps an unexpected internal error.
// The original error is hidden from clients but preserved for logging.
func Internal(cause error) *AppError {
	return New(http.StatusInternalServerError, "an unexpected error occurred", cause)
}

// Wrap is like Internal but lets you customize the message.
func Wrap(code int, msg string, cause error) *AppError {
	return New(code, msg, cause)
}

// ---------------------------------------------------------------------------
// Inspection helpers
// ---------------------------------------------------------------------------

// IsNotFound reports whether err is a 404 AppError.
func IsNotFound(err error) bool {
	var ae *AppError
	return errors.As(err, &ae) && ae.Code == http.StatusNotFound
}

// IsConflict reports whether err is a 409 AppError.
func IsConflict(err error) bool {
	var ae *AppError
	return errors.As(err, &ae) && ae.Code == http.StatusConflict
}

// IsUnauthorized reports whether err is a 401 AppError.
func IsUnauthorized(err error) bool {
	var ae *AppError
	return errors.As(err, &ae) && ae.Code == http.StatusUnauthorized
}

// As extracts the *AppError from an error chain (convenience wrapper).
func As(err error) (*AppError, bool) {
	var ae *AppError
	ok := errors.As(err, &ae)
	return ae, ok
}
