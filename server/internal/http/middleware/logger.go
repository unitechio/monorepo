package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = time.Now().Format("20060102150405")
		}
		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-ID", correlationID)

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		requestID, _ := c.Get("RequestID")

		attrs := []any{
			slog.String("correlation_id", correlationID),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
			slog.String("request_id", toString(requestID)),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("gin_errors", c.Errors.String()))
		}

		switch {
		case statusCode >= 500:
			logger.Error("request", attrs...)
		case statusCode >= 400:
			logger.Warn("request", attrs...)
		default:
			logger.Info("request", attrs...)
		}
	}
}
