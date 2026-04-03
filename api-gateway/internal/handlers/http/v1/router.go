package v1

import (
	"fmt"
	"time"

	"github.com/gfdmit/web-forum/api-gateway/config"
	"github.com/gfdmit/web-forum/api-gateway/internal/middleware"
	"github.com/gfdmit/web-forum/api-gateway/internal/proxy"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New(conf *config.Config) (*gin.Engine, error) {
	var (
		router = gin.New()
	)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           2 * time.Hour,
	}))

	api := router.Group("/api/v1")
	{
		api.Use(gin.Logger())

		postGroup := api.Group("")
		{
			postProxy, err := proxy.New(fmt.Sprintf("http://%s:%s", conf.PostService.Host, conf.PostService.Port))
			if err != nil {
				return nil, fmt.Errorf("error when setup router: %v", err)
			}
			postGroup.Use(middleware.RewritePrefix("/api/v1", "/api/v2"))

			postGroup.GET("/boards", postProxy.Forward())
			postGroup.GET("/boards/:id", postProxy.Forward())
			postGroup.POST("/boards", postProxy.Forward())
			postGroup.DELETE("/boards/:id", postProxy.Forward())
			postGroup.POST("/boards/:id/restore", postProxy.Forward())

			postGroup.GET("/boards/:id/posts", postProxy.Forward())
			postGroup.GET("/posts/:id", postProxy.Forward())
			postGroup.POST("/posts", postProxy.Forward())
			postGroup.DELETE("/posts/:id", postProxy.Forward())

			postGroup.GET("/posts/:id/comments", postProxy.Forward())
			postGroup.GET("/comments/:id", postProxy.Forward())
			postGroup.POST("/comments", postProxy.Forward())
			postGroup.DELETE("/comments/:id", postProxy.Forward())

			postGroup.GET("/ping", postProxy.Forward())
		}

	}

	return router, nil
}
