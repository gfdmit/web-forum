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
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           2 * time.Hour,
	}))

	postProxy, err := proxy.New(fmt.Sprintf("http://%s:%s", conf.PostService.Host, conf.PostService.Port))
	if err != nil {
		return nil, fmt.Errorf("error when setup router: %v", err)
	}

	authProxy, err := proxy.New(fmt.Sprintf("http://%s:%s", conf.AuthService.Host, conf.AuthService.Port))
	if err != nil {
		return nil, fmt.Errorf("error when setup router: %v", err)
	}

	api := router.Group("/api/v1")
	api.Use(gin.Logger())

	public := api.Group("")
	{
		public.POST("/auth/login", authProxy.Forward())

		public.Use(middleware.RewritePrefix("/api/v1", "/api/v2"))

		public.GET("/ping", postProxy.Forward())

		public.GET("/boards", postProxy.Forward())
		public.GET("/boards/:id", postProxy.Forward())

		public.GET("/boards/:id/posts", postProxy.Forward())
		public.GET("/posts/:id", postProxy.Forward())

		public.GET("/posts/:id/comments", postProxy.Forward())
		public.GET("/comments/:id", postProxy.Forward())
	}

	protected := api.Group("")
	{
		protected.Use(middleware.RewritePrefix("/api/v1", "/api/v2"))
		protected.Use(middleware.Auth(conf.JWT.Secret))

		protected.POST("/boards", postProxy.Forward())
		protected.DELETE("/boards/:id", postProxy.Forward())
		protected.POST("/boards/:id/restore", postProxy.Forward())

		protected.POST("/posts", postProxy.Forward())
		protected.DELETE("/posts/:id", postProxy.Forward())

		protected.POST("/comments", postProxy.Forward())
		protected.DELETE("/comments/:id", postProxy.Forward())
	}

	return router, nil
}
