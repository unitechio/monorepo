// Package response provides a standardized JSON envelope for all API responses.
//
// Every response — success or error — has this shape:
//
//	{
//	  "success": true|false,
//	  "message": "Human-readable message",
//	  "data":    <payload or null>,
//	  "error":   <error detail or null>,
//	  "meta":    <pagination meta or null>,
//	  "request_id": "uuid"
//	}
//
// This makes it easy for the frontend to check response.success and always
// find the payload in response.data.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// Envelope types
// ---------------------------------------------------------------------------

// Envelope is the standard JSON wrapper for every API response.
type Envelope struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	Meta      *Meta      `json:"meta,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
}

// ErrorBody carries structured error details (never internal stack traces).
type ErrorBody struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"` // validation field errors
}

// Meta holds pagination information for list endpoints.
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// ---------------------------------------------------------------------------
// Success helpers
// ---------------------------------------------------------------------------

// OK writes a 200 JSON success envelope.
func OK(c *gin.Context, message string, data any) {
	respond(c, http.StatusOK, message, data, nil, nil)
}

// Created writes a 201 JSON success envelope.
func Created(c *gin.Context, message string, data any) {
	respond(c, http.StatusCreated, message, data, nil, nil)
}

// NoContent writes a 204 (no body).
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// OKWithMeta writes a 200 JSON success envelope including pagination metadata.
func OKWithMeta(c *gin.Context, message string, data any, meta *Meta) {
	respond(c, http.StatusOK, message, data, nil, meta)
}

// ---------------------------------------------------------------------------
// Error helpers
// ---------------------------------------------------------------------------

// Fail writes an error envelope with the provided HTTP status and message.
func Fail(c *gin.Context, code int, message string) {
	respond(c, code, message, nil, &ErrorBody{Code: code, Message: message}, nil)
}

// FailWithFields writes a 400 error envelope with per-field validation messages.
func FailWithFields(c *gin.Context, message string, fields map[string]string) {
	respond(c, http.StatusBadRequest, message, nil, &ErrorBody{
		Code:    http.StatusBadRequest,
		Message: message,
		Fields:  fields,
	}, nil)
}

// InternalError writes a 500 error. The internal cause is logged but never
// exposed to the client.
func InternalError(c *gin.Context) {
	Fail(c, http.StatusInternalServerError, "an unexpected error occurred")
}

// ---------------------------------------------------------------------------
// Core dispatcher
// ---------------------------------------------------------------------------

// respond is the single write point for all JSON responses.
// It injects the request_id from the gin context (set by RequestIDMiddleware).
func respond(c *gin.Context, code int, message string, data any, errBody *ErrorBody, meta *Meta) {
	requestID, _ := c.Get("RequestID")

	payload := Envelope{
		Success:   errBody == nil,
		Message:   message,
		Data:      data,
		Error:     errBody,
		Meta:      meta,
		RequestID: toString(requestID),
	}

	c.JSON(code, payload)
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
