package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gfdmit/web-forum/post-service/internal/model"
	"github.com/gfdmit/web-forum/post-service/internal/repository"
)

var (
	ErrValidation = errors.New("validation error")
)

type Service interface {
	GetBoard(ctx context.Context, id int) (model.Board, error)
	GetBoards(ctx context.Context, includeDeleted bool) ([]model.Board, error)
	CreateBoard(ctx context.Context, input model.CreateBoardInput) (model.Board, error)
	DeleteBoard(ctx context.Context, id int) error
	RestoreBoard(ctx context.Context, id int) error

	GetPost(ctx context.Context, id int) (model.Post, error)
	GetPosts(ctx context.Context, boardID int, includeDeleted bool, limit, offset int) ([]model.Post, error)
	CreatePost(ctx context.Context, input model.CreatePostInput) (model.Post, error)
	DeletePost(ctx context.Context, id int) error

	GetComment(ctx context.Context, id int) (model.Comment, error)
	GetComments(ctx context.Context, postID int, includeDeleted bool, limit, offset int) ([]model.Comment, error)
	CreateComment(ctx context.Context, input model.CreateCommentInput) (model.Comment, error)
	DeleteComment(ctx context.Context, id int) error

	GetProfile(ctx context.Context, userID int) (model.Profile, error)
	GetProfiles(ctx context.Context, includeDeleted bool) ([]model.Profile, error)
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (svc *service) GetBoard(ctx context.Context, id int) (model.Board, error) {
	return svc.repo.GetBoard(ctx, id)
}

func (svc *service) GetBoards(ctx context.Context, includeDeleted bool) ([]model.Board, error) {
	return svc.repo.GetBoards(ctx, includeDeleted)
}

func (svc *service) CreateBoard(ctx context.Context, input model.CreateBoardInput) (model.Board, error) {
	if input.Name == "" {
		return model.Board{}, ErrValidation
	}
	return svc.repo.CreateBoard(ctx, input)
}

func (svc *service) DeleteBoard(ctx context.Context, id int) error {
	return svc.repo.DeleteBoard(ctx, id)
}

func (svc *service) RestoreBoard(ctx context.Context, id int) error {
	return svc.repo.RestoreBoard(ctx, id)
}

func (svc *service) GetPost(ctx context.Context, id int) (model.Post, error) {
	return svc.repo.GetPost(ctx, id)
}

func (svc *service) GetPosts(ctx context.Context, boardID int, includeDeleted bool, limit, offset int) ([]model.Post, error) {
	if limit > 100 || limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return svc.repo.GetPosts(ctx, boardID, includeDeleted, limit, offset)
}

func (svc *service) CreatePost(ctx context.Context, input model.CreatePostInput) (model.Post, error) {
	if len(input.Title) > 100 {
		return model.Post{}, ErrValidation
	}
	if len(input.Text) > 5000 {
		return model.Post{}, ErrValidation
	}
	if _, err := svc.repo.GetBoard(ctx, input.BoardID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Post{}, repository.ErrNotFound
		}
		return model.Post{}, fmt.Errorf("CreatePost check board: %w", err)
	}
	return svc.repo.CreatePost(ctx, input)
}

func (svc *service) DeletePost(ctx context.Context, id int) error {
	return svc.repo.DeletePost(ctx, id)
}

func (svc *service) GetComment(ctx context.Context, id int) (model.Comment, error) {
	return svc.repo.GetComment(ctx, id)
}

func (svc *service) GetComments(ctx context.Context, postID int, includeDeleted bool, limit, offset int) ([]model.Comment, error) {
	if limit > 100 || limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return svc.repo.GetComments(ctx, postID, includeDeleted, limit, offset)
}

func (svc *service) CreateComment(ctx context.Context, input model.CreateCommentInput) (model.Comment, error) {
	if len(input.Text) > 5000 {
		return model.Comment{}, ErrValidation
	}
	if _, err := svc.repo.GetPost(ctx, input.PostID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Comment{}, repository.ErrNotFound
		}
		return model.Comment{}, fmt.Errorf("CreateComment check post: %w", err)
	}
	return svc.repo.CreateComment(ctx, input)
}

func (svc *service) DeleteComment(ctx context.Context, id int) error {
	return svc.repo.DeleteComment(ctx, id)
}

func (svc *service) GetProfile(ctx context.Context, userID int) (model.Profile, error) {
	return svc.repo.GetProfile(ctx, userID)
}

func (svc *service) GetProfiles(ctx context.Context, includeDeleted bool) ([]model.Profile, error) {
	return svc.repo.GetProfiles(ctx, includeDeleted)
}
