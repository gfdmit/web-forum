package rest

import (
	"net/http"

	"github.com/gfdmit/web-forum/post-service/internal/service"
	"github.com/gin-gonic/gin"
)

type restHandler struct {
	svc *service.Service
}

func New(svc *service.Service) (*restHandler, error) {
	return &restHandler{
		svc: svc,
	}, nil
}

func (rh *restHandler) GetBoards(c *gin.Context) {
	includeDeleted := c.Query("includeDeleted") == "true"

	boards, err := rh.svc.GetBoards(c.Request.Context(), includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, boards)
}

func (rh *restHandler) GetBoard(c *gin.Context) {
	id := c.Param("id")

	board, err := rh.svc.GetBoard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, board)
}

func (rh *restHandler) CreateBoard(c *gin.Context) {
	var body struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	board, err := rh.svc.CreateBoard(c.Request.Context(), body.Name, body.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, board)
}

func (rh *restHandler) DeleteBoard(c *gin.Context) {
	id := c.Param("id")

	ok, err := rh.svc.DeleteBoard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": ok})
}

func (rh *restHandler) RestoreBoard(c *gin.Context) {
	id := c.Param("id")

	ok, err := rh.svc.RestoreBoard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": ok})
}

func (rh *restHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := rh.svc.GetPost(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, post)
}

func (rh *restHandler) GetPosts(c *gin.Context) {
	boardId := c.Param("id")

	includeDeleted := c.Query("includeDeleted") == "true"
	limit := queryInt(c.Query("limit"), 100)
	offset := queryInt(c.Query("offset"), 0)

	posts, err := rh.svc.GetPosts(c.Request.Context(), boardId, includeDeleted, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (rh *restHandler) CreatePost(c *gin.Context) {
	var body struct {
		BoardId string `json:"boardId" binding:"required"`
		Title   string `json:"title"`
		Text    string `json:"text" binding:"required"`
		HashIp  string `json:"hashIp"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := rh.svc.CreatePost(c.Request.Context(), body.BoardId, body.Title, body.Text, body.HashIp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (rh *restHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	ok, err := rh.svc.DeletePost(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": ok})
}

func (rh *restHandler) GetComment(c *gin.Context) {
	id := c.Param("id")

	comment, err := rh.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (rh *restHandler) GetComments(c *gin.Context) {
	postId := c.Param("id")

	includeDeleted := c.Query("includeDeleted") == "true"
	limit := queryInt(c.Query("limit"), 500)
	offset := queryInt(c.Query("offset"), 0)

	comments, err := rh.svc.GetComments(c.Request.Context(), postId, includeDeleted, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (rh *restHandler) CreateComment(c *gin.Context) {
	var body struct {
		PostId string `json:"postId" binding:"required"`
		Text   string `json:"text" binding:"required"`
		HashIp string `json:"hashIp"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	comment, err := rh.svc.CreateComment(c.Request.Context(), body.PostId, body.Text, body.HashIp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (rh *restHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	ok, err := rh.svc.DeleteComment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": ok})
}
