package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gfdmit/web-forum/post-service/internal/repository"
	"github.com/gfdmit/web-forum/post-service/internal/service"
	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse(err))
	case errors.Is(err, service.ErrValidation):
		c.JSON(http.StatusBadRequest, errorResponse(err))
	default:
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
}

func parseID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func parseUserID(c *gin.Context) (int, error) {
	val := c.GetHeader("X-User-Id")
	if val == "" {
		return 0, errors.New("missing X-User-Id header")
	}
	id, err := strconv.Atoi(val)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid X-User-Id header")
	}
	return id, nil
}

func queryInt(val string, def int) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return n
}
