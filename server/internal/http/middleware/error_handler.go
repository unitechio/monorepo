package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		requestID, _ := c.Get(response.ContextKeyRequestID)
		rid := toString(requestID)
		appErr, ok := apperr.As(err)
		if !ok {
			appErr = apperr.Internal(err)
		}

		logger.Error("request failed",
			slog.String("request_id", rid),
			slog.Int("status", appErr.Code),
			slog.String("message", appErr.Message),
			slog.String("error", err.Error()),
		)
		response.FailError(c, appErr)
	}
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
