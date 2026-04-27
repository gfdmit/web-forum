package handler

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
	router.Use(gin.Logger(), gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           2 * time.Hour,
	}))

	router.GET("/ping", pingHandler(map[string]string{
		"auth":  fmt.Sprintf("http://%s:%s/ping", conf.AuthService.Host, conf.AuthService.Port),
		"posts": fmt.Sprintf("http://%s:%s/ping", conf.PostService.Host, conf.PostService.Port),
	}))

	postProxy, err := proxy.New(fmt.Sprintf("http://%s:%s", conf.PostService.Host, conf.PostService.Port))
	if err != nil {
		return nil, fmt.Errorf("proxy.New post: %w", err)
	}

	authProxy, err := proxy.New(fmt.Sprintf("http://%s:%s", conf.AuthService.Host, conf.AuthService.Port))
	if err != nil {
		return nil, fmt.Errorf("proxy.New auth: %w", err)
	}

	api := router.Group("/api/v1")

	public := api.Group("")
	{
		public.POST("/auth/login", authProxy.Forward())

		public.Use(middleware.RewritePrefix("/api/v1", "/api/v2"))

		public.GET("/boards", postProxy.Forward())
		public.GET("/boards/:id", postProxy.Forward())

		public.GET("/boards/:id/posts", postProxy.Forward())
		public.GET("/posts/:id", postProxy.Forward())

		public.GET("/posts/:id/comments", postProxy.Forward())
		public.GET("/comments/:id", postProxy.Forward())

		public.GET("profiles", postProxy.Forward())
		public.GET("profiles/:id", postProxy.Forward())
	}

	protected := api.Group("")
	protected.Use(middleware.Auth(conf.JWT.Secret))
	{
		public.POST("/auth/logout", authProxy.Forward())

		protected.Use(middleware.RewritePrefix("/api/v1", "/api/v2"))

		protected.POST("/boards", postProxy.Forward())
		protected.DELETE("/boards/:id", postProxy.Forward())
		protected.POST("/boards/:id/restore", postProxy.Forward())

		protected.POST("/posts", postProxy.Forward())
		protected.DELETE("/posts/:id", postProxy.Forward())

		protected.POST("/comments", postProxy.Forward())
		protected.DELETE("/comments/:id", postProxy.Forward())

		protected.GET("/profile", postProxy.Forward())
	}

	return router, nil
}
