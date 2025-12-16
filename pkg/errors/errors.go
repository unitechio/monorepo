package errors

import (
	"errors"
)

// Common errors
var (
	ErrUnauthorized            = errors.New("unauthorized")
	ErrInvalidToken            = errors.New("invalid token")
	ErrForbidden               = errors.New("forbidden")
	ErrNotFound                = errors.New("not found")
	ErrBadRequest              = errors.New("bad request")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrInternal                = errors.New("internal error")
)
