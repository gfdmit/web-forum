package rest

import (
	"net/http"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type restHandler struct {
	conf *config.JWT
	svc  *service.Service
}

func New(conf *config.JWT, svc *service.Service) *restHandler {
	return &restHandler{
		conf: conf,
		svc:  svc,
	}
}

func (rh *restHandler) Login(c *gin.Context) {
	type loginRequest struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := rh.svc.GenerateToken(req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.SetCookie("token", token, int(rh.conf.TTL.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
