package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

type visitor struct {
	count   int
	resetAt time.Time
}

func RateLimit(rps int) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	return func(c *gin.Context) {
		now := time.Now()
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists || now.After(v.resetAt) {
			visitors[ip] = &visitor{count: 1, resetAt: now.Add(time.Minute)}
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rps))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rps-1))
			mu.Unlock()
			c.Next()
			return
		}
		v.count++
		remaining := rps - v.count
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rps))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", v.resetAt.Unix()))
		if v.count > rps {
			mu.Unlock()
			response.FailError(c, apperr.TooManyRequests("too many requests"))
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
}
