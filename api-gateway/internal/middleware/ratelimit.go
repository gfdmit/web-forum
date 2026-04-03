// internal/middleware/ratelimit.go
package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters sync.Map
)

func getLimiter(ip string, rps int) *rate.Limiter {
	if l, ok := limiters.Load(ip); ok {
		return l.(*rate.Limiter)
	}
	l := rate.NewLimiter(rate.Limit(rps), rps)
	limiters.Store(ip, l)
	return l
}

func RateLimit(rps int) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getLimiter(c.ClientIP(), rps)
		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
