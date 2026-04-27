package rest

import (
	"net/http"

	"github.com/gfdmit/web-forum/post-service/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(svc service.Service) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	h := New(svc)

	api := router.Group("/api/v2")
	{
		api.GET("/boards", h.GetBoards)
		api.GET("/boards/:id", h.GetBoard)
		api.POST("/boards", h.CreateBoard)
		api.DELETE("/boards/:id", h.DeleteBoard)
		api.POST("/boards/:id/restore", h.RestoreBoard)

		api.GET("/boards/:id/posts", h.GetPosts)
		api.GET("/posts/:id", h.GetPost)
		api.POST("/posts", h.CreatePost)
		api.DELETE("/posts/:id", h.DeletePost)

		api.GET("/posts/:id/comments", h.GetComments)
		api.GET("/comments/:id", h.GetComment)
		api.POST("/comments", h.CreateComment)
		api.DELETE("/comments/:id", h.DeleteComment)

		api.GET("/profile", h.GetMyProfile)
		api.GET("/profiles", h.GetProfiles)
		api.GET("/profiles/:id", h.GetProfile)
	}

	return router
}
