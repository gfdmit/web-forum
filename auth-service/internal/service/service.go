package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/kfu"
	"github.com/gfdmit/web-forum/auth-service/internal/model"
	"github.com/gfdmit/web-forum/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	GenerateToken(ctx context.Context, login, password string) (string, error)
}

type service struct {
	conf *config.JWT
	repo repository.Repository
	kfu  kfu.KFUClient
}

func New(conf *config.JWT, repo repository.Repository, kfu kfu.KFUClient) Service {
	return &service{conf: conf, repo: repo, kfu: kfu}
}

func (svc *service) GenerateToken(ctx context.Context, login, password string) (string, error) {
	id, err := svc.resolveUser(ctx, login, password)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": strconv.Itoa(id),
		"exp": time.Now().Add(svc.conf.TTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(svc.conf.Secret))
	if err != nil {
		return "", fmt.Errorf("GenerateToken sign: %w", err)
	}
	return signed, nil
}

func (svc *service) resolveUser(ctx context.Context, login, password string) (int, error) {
	user, err := svc.repo.GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return 0, fmt.Errorf("resolveUser get: %w", err)
	}

	needsUpdate := errors.Is(err, repository.ErrNotFound) ||
		bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil

	if !needsUpdate {
		return user.ID, nil
	}

	return svc.syncFromKFU(ctx, login, password)
}

func (svc *service) syncFromKFU(ctx context.Context, login, password string) (int, error) {
	profile, err := svc.kfu.ParseKFU(ctx, login, password)
	if err != nil {
		if errors.Is(err, kfu.ErrInvalidCredentials) {
			return 0, ErrInvalidCredentials
		}
		return 0, fmt.Errorf("syncFromKFU parse: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("syncFromKFU hash: %w", err)
	}

	user, err := svc.repo.CreateOrUpdateUser(ctx, model.CreateUserInput{
		Login:        login,
		PasswordHash: string(hash),
	})
	if err != nil {
		return 0, fmt.Errorf("syncFromKFU user: %w", err)
	}

	if err := svc.repo.CreateOrUpdateProfile(ctx, toProfileInput(user.ID, profile)); err != nil {
		return 0, fmt.Errorf("syncFromKFU profile: %w", err)
	}

	return user.ID, nil
}

func toProfileInput(userID int, p *model.Profile) model.CreateProfileInput {
	return model.CreateProfileInput{
		UserID:       userID,
		UniversityID: p.AllId,
		Firstname:    p.Firstname,
		Lastname:     p.Lastname,
		Middlename:   p.MiddleName,
		Birthday:     p.Birthday,
		Faculty:      p.Faculty,
		Grade:        p.Grade,
		Group:        p.Group,
		Status:       p.Status,
	}
}
