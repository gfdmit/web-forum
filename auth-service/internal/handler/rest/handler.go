package rest

import (
	"net/http"
	"time"

	"github.com/gfdmit/web-forum/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Handler struct {
	svc service.Service
	ttl time.Duration
}

func New(svc service.Service, ttl time.Duration) *Handler {
	return &Handler{svc: svc, ttl: ttl}
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.GenerateToken(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	c.SetCookie("token", token, int(h.ttl.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
