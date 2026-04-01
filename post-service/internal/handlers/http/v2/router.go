package v2

import (
	"net/http"
	"time"

	"github.com/gfdmit/web-forum/post-service/internal/handlers/http/v2/rest"
	"github.com/gfdmit/web-forum/post-service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New(svc *service.Service) (*gin.Engine, error) {
	var (
		router = gin.New()
	)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300 * time.Second,
	}))

	restHandler, err := rest.New(svc)
	if err != nil {
		return nil, err
	}

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.Use(gin.Logger())

		authGroup := apiGroup.Group("")
		{
			apiGroup.GET("/boards", restHandler.GetBoards)
			apiGroup.GET("/boards/:id", restHandler.GetBoard)
			apiGroup.POST("/boards", restHandler.CreateBoard)
			apiGroup.DELETE("/boards/:id", restHandler.DeleteBoard)
			apiGroup.POST("/boards/:id/restore", restHandler.RestoreBoard)
		}

		authGroup = apiGroup.Group("")
		{
			apiGroup.GET("/boards/:id/posts", restHandler.GetPosts)
			apiGroup.GET("/posts/:id", restHandler.GetPost)
			apiGroup.POST("/posts", restHandler.CreatePost)
			apiGroup.DELETE("/posts/:id", restHandler.DeletePost)
		}

		authGroup = apiGroup.Group("")
		{
			apiGroup.GET("/posts/:id/comments", restHandler.GetComments)
			apiGroup.GET("/comments/:id", restHandler.GetComment)
			apiGroup.POST("/comments", restHandler.CreateComment)
			apiGroup.DELETE("/comments/:id", restHandler.DeleteComment)
		}

		authGroup = apiGroup.Group("")
		{
			authGroup.GET("/ping", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
		}
	}

	return router, nil
}
