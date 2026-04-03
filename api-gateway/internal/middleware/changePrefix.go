package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func RewritePrefix(from, to string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = to + strings.TrimPrefix(c.Request.URL.Path, from)
		c.Next()
	}
}
