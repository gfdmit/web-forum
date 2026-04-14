package rest

import (
	"net/http"
	"time"

	"github.com/gfdmit/web-forum/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(svc service.Service, ttl time.Duration) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	h := New(svc, ttl)

	api := router.Group("/api/v1")
	{
		api.POST("/auth/login", h.Login)
	}

	return router
}
