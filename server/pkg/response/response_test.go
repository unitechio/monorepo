package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

func TestFailError_AppErrorUsesStandardEnvelope(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set(response.ContextKeyRequestID, "req-123")

	response.FailError(c, apperr.Unauthorized("invalid credentials"))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var payload response.Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if payload.Success {
		t.Fatal("expected success=false")
	}
	if payload.Message != "invalid credentials" {
		t.Fatalf("expected message %q, got %q", "invalid credentials", payload.Message)
	}
	if payload.Error == nil {
		t.Fatal("expected error body")
	}
	if payload.Error.Code != http.StatusUnauthorized {
		t.Fatalf("expected error code %d, got %d", http.StatusUnauthorized, payload.Error.Code)
	}
	if payload.RequestID != "req-123" {
		t.Fatalf("expected request id %q, got %q", "req-123", payload.RequestID)
	}
}

func TestFailError_UnknownErrorBecomesInternalError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	response.FailError(c, errors.New("database exploded"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var payload response.Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if payload.Error == nil {
		t.Fatal("expected error body")
	}
	if payload.Error.Message != "an unexpected error occurred" {
		t.Fatalf("expected sanitized message, got %q", payload.Error.Message)
	}
}
