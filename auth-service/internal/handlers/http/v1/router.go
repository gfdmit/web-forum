package v1

import (
	"net/http"
	"time"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/handlers/http/v1/rest"
	"github.com/gfdmit/web-forum/auth-service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New(conf *config.JWT, svc *service.Service) (*gin.Engine, error) {
	var (
		router = gin.New()
	)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           2 * time.Hour,
	}))

	restHandler := rest.New(conf, svc)

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.Use(gin.Logger())

		apiGroup.POST("/auth/login", restHandler.Login)

		apiGroup.GET("/ping", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	}

	return router, nil
}
