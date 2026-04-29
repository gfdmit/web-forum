package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	maxFileSize = 10 << 20
)

var allowedTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, errorResponse(err))
}

func extensionToContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
