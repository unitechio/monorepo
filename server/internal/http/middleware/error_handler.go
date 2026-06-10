package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		requestID, _ := c.Get("RequestID")
		rid := toString(requestID)

		var ae *apperr.AppError
		if errors.As(err, &ae) {
			if ae.Code >= 500 {
				logger.Error("application error",
					slog.String("request_id", rid),
					slog.String("error", ae.Error()),
				)
			}
			response.Fail(c, ae.Code, ae.Message)
			return
		}

		logger.Error("unhandled error",
			slog.String("request_id", rid),
			slog.String("error", err.Error()),
		)
		response.Fail(c, http.StatusInternalServerError, "an unexpected error occurred")
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
