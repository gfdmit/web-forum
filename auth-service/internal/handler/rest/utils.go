package rest

import (
	"errors"
	"net/http"

	"github.com/gfdmit/web-forum/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, errorResponse(err))
	default:
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
}
