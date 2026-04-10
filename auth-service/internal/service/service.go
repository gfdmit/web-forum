package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/repository/postgres"
	"github.com/gfdmit/web-forum/auth-service/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	id, passHash, err := svc.repo.GetIdPassHash(login)
	if err != nil && !errors.Is(err, postgres.ErrNotFound) {
		fmt.Printf("db dead: %v\n", err)
		return "", err
	}

	isParser := errors.Is(err, postgres.ErrNotFound) || (bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password)) != nil)
	if isParser {
		profile, err := utils.ParseKFU(login, password)
		if err != nil {
			fmt.Printf("parser dead: %v\n", err)
			return "", err
		}
		newPassHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("hasher dead: %v\n", err)
			return "", err
		}
		id, err = svc.repo.CreateOrUpdateUser(login, string(newPassHash))
		if err != nil {
			fmt.Printf("users dead: %v\n", err)
			return "", err
		}
		if err := svc.repo.CreateOrUpdateProfile(id, profile); err != nil {
			fmt.Printf("profile dead: %v\n", err)
			return "", err
		}
	}

	claims := jwt.MapClaims{
		"sub": strconv.Itoa(id),
		"exp": time.Now().Add(svc.conf.TTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(svc.conf.Secret))
}
