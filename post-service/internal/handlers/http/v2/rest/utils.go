package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func queryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
		return defaultVal
	}
	return n
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
