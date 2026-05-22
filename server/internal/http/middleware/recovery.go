package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get(response.ContextKeyRequestID)
				rid := toString(requestID)
				appErr := apperr.Internal(fmt.Errorf("panic: %v", err))

				logger.Error("panic recovered",
					slog.String("request_id", rid),
					slog.Int("status", appErr.Code),
					slog.String("message", appErr.Message),
					slog.Any("error", err),
				)

				response.FailError(c, appErr)
				c.Abort()
			}
		}()
		c.Next()
	}
}
