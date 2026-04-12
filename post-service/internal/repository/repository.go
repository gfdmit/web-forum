package repository

import (
	"context"
	"errors"

	"github.com/gfdmit/web-forum/post-service/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
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
}
