package service

import (
	"errors"
	"time"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	conf *config.JWT
	repo *postgres.Repository
}

func New(conf *config.JWT, repo *postgres.Repository) *Service {
	return &Service{
		conf: conf,
		repo: repo,
	}
}

func (svc *Service) GenerateToken(login, password string) (string, error) {
	if login != "admin" || password != "admin" {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": login,
		"exp": time.Now().Add(svc.conf.TTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(svc.conf.Secret))
}
