package rest

import (
	"net/http"

	"github.com/gfdmit/web-forum/post-service/internal/model"
	"github.com/gfdmit/web-forum/post-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetBoards(c *gin.Context) {
	includeDeleted := c.Query("includeDeleted") == "true"

	boards, err := h.svc.GetBoards(c.Request.Context(), includeDeleted)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, boards)
}

func (h *Handler) GetBoard(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	board, err := h.svc.GetBoard(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, board)
}

func (h *Handler) CreateBoard(c *gin.Context) {
	var input model.CreateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	board, err := h.svc.CreateBoard(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, board)
}

func (h *Handler) DeleteBoard(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.svc.DeleteBoard(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) RestoreBoard(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.svc.RestoreBoard(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetPosts(c *gin.Context) {
	boardID, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	includeDeleted := c.Query("includeDeleted") == "true"
	limit := queryInt(c.Query("limit"), 20)
	offset := queryInt(c.Query("offset"), 0)

	posts, err := h.svc.GetPosts(c.Request.Context(), boardID, includeDeleted, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPost(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := h.svc.GetPost(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *Handler) CreatePost(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var input model.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	input.UserID = &userID

	post, err := h.svc.CreatePost(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *Handler) DeletePost(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.svc.DeletePost(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetComments(c *gin.Context) {
	postID, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	includeDeleted := c.Query("includeDeleted") == "true"
	limit := queryInt(c.Query("limit"), 20)
	offset := queryInt(c.Query("offset"), 0)

	comments, err := h.svc.GetComments(c.Request.Context(), postID, includeDeleted, limit, offset)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *Handler) GetComment(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	comment, err := h.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (h *Handler) CreateComment(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var input model.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	input.UserID = &userID

	comment, err := h.svc.CreateComment(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.svc.DeleteComment(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetProfiles(c *gin.Context) {
	includeDeleted := c.Query("includeDeleted") == "true"

	profiles, err := h.svc.GetProfiles(c.Request.Context(), includeDeleted)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *Handler) GetProfile(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	profile, err := h.svc.GetProfile(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *Handler) GetMyProfile(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	profile, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}
