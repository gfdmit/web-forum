package repository

import (
	"context"
	"errors"

	"github.com/gfdmit/web-forum/auth-service/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	CreateOrUpdateUser(ctx context.Context, input model.CreateUserInput) (model.User, error)
	CreateOrUpdateProfile(ctx context.Context, input model.CreateProfileInput) error
}
