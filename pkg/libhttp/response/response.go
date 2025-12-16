package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success creates a successful response
func Success(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

// Error creates an error response
func Error(code, message string, details interface{}) *Response {
	return &Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// Error writes an error response to gin context
func ErrorResponse(c *gin.Context, err error) {
	var code string
	var statusCode int

	switch err.Error() {
	case "unauthorized":
		code = "UNAUTHORIZED"
		statusCode = http.StatusUnauthorized
	case "invalid token":
		code = "INVALID_TOKEN"
		statusCode = http.StatusUnauthorized
	case "forbidden":
		code = "FORBIDDEN"
		statusCode = http.StatusForbidden
	case "not found":
		code = "NOT_FOUND"
		statusCode = http.StatusNotFound
	case "bad request":
		code = "BAD_REQUEST"
		statusCode = http.StatusBadRequest
	case "insufficient permissions":
		code = "INSUFFICIENT_PERMISSIONS"
		statusCode = http.StatusForbidden
	case "internal error":
		code = "INTERNAL_ERROR"
		statusCode = http.StatusInternalServerError
	default:
		code = "ERROR"
		statusCode = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: err.Error(),
		},
	})
}
