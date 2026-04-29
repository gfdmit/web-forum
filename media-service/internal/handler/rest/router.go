package rest

import (
	"net/http"

	"github.com/gfdmit/web-forum/media-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func NewRouter(store storage.Storage, publicHost string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	h := New(store, publicHost)

	api := router.Group("/api/v1")
	{
		api.POST("/media/upload", h.PostMedia)
		api.GET("/media/:filename", h.GetMedia)
	}

	return router
}
