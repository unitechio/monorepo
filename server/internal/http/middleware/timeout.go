package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.Abort()
			response.FailError(c, apperr.Wrap(http.StatusGatewayTimeout, "request timeout", ctx.Err()))
		}
	}
}
