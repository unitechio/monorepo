package apperr_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
)

func TestAppError_Error(t *testing.T) {
	err := apperr.New(400, "bad input")
	if err.Error() != "[400] bad input" {
		t.Errorf("unexpected: %s", err.Error())
	}
}

func TestBadRequest(t *testing.T) {
	err := apperr.BadRequest("invalid email")
	if err.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", err.Code)
	}
	if err.Message != "invalid email" {
		t.Errorf("unexpected message: %s", err.Message)
	}
}

func TestNotFound(t *testing.T) {
	err := apperr.NotFound("user", 42)
	if err.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", err.Code)
	}
}

func TestConflict(t *testing.T) {
	err := apperr.Conflict("email taken")
	if !apperr.IsConflict(err) {
		t.Error("expected IsConflict to be true")
	}
}

func TestUnauthorized(t *testing.T) {
	err := apperr.Unauthorized("no token")
	if !apperr.IsUnauthorized(err) {
		t.Error("expected IsUnauthorized to be true")
	}
}

func TestInternal_HidesMessage(t *testing.T) {
	cause := errors.New("db connection refused")
	err := apperr.Internal(cause)
	if err.Message != "an unexpected error occurred" {
		t.Errorf("internal cause leaked: %s", err.Message)
	}
	if err.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", err.Code)
	}
}

func TestUnwrap(t *testing.T) {
	root := errors.New("root cause")
	wrapped := apperr.Internal(root)
	if !errors.Is(wrapped, root) {
		t.Error("expected errors.Is to unwrap AppError to root cause")
	}
}

func TestAs(t *testing.T) {
	err := apperr.BadRequest("oops")
	ae, ok := apperr.As(err)
	if !ok {
		t.Fatal("expected ok")
	}
	if ae.Code != 400 {
		t.Errorf("expected 400, got %d", ae.Code)
	}
}
