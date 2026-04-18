package service

import (
	"strings"
	"testing"

	"github.com/gfdmit/web-forum/post-service/internal/mocks"
	"github.com/gfdmit/web-forum/post-service/internal/model"
	"github.com/gfdmit/web-forum/post-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ───── CreateBoard ─────

func TestCreateBoard_EmptyName(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	_, err := svc.CreateBoard(t.Context(), model.CreateBoardInput{Name: ""})

	assert.ErrorIs(t, err, ErrValidation)
}

func TestCreateBoard_Success(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateBoardInput{Name: "General"}
	expected := model.Board{ID: 1, Name: "General"}

	repo.EXPECT().CreateBoard(t.Context(), input).Return(expected, nil)

	board, err := svc.CreateBoard(t.Context(), input)

	require.NoError(t, err)
	assert.Equal(t, expected, board)
}

// ───── CreatePost ─────

func TestCreatePost_TitleTooLong(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{
		Title:   strings.Repeat("a", 101),
		Text:    "text",
		BoardID: 1,
	}

	_, err := svc.CreatePost(t.Context(), input)

	assert.ErrorIs(t, err, ErrValidation)
}

func TestCreatePost_TextTooLong(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{
		Title:   "Normal title",
		Text:    strings.Repeat("x", 5001),
		BoardID: 1,
	}

	_, err := svc.CreatePost(t.Context(), input)

	assert.ErrorIs(t, err, ErrValidation)
}

func TestCreatePost_BoardNotFound(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{Text: "text", BoardID: 99}
	repo.EXPECT().GetBoard(t.Context(), 99).Return(model.Board{}, repository.ErrNotFound)

	_, err := svc.CreatePost(t.Context(), input)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestCreatePost_RepoError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{Text: "text", BoardID: 1}
	repo.EXPECT().GetBoard(t.Context(), 1).Return(model.Board{}, assert.AnError)

	_, err := svc.CreatePost(t.Context(), input)

	assert.ErrorIs(t, err, assert.AnError)
}

func TestCreatePost_Success(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{Text: "text", BoardID: 1}
	expected := model.Post{ID: 1, Text: "text", BoardID: 1}

	repo.EXPECT().GetBoard(t.Context(), 1).Return(model.Board{ID: 1}, nil)
	repo.EXPECT().CreatePost(t.Context(), input).Return(expected, nil)

	post, err := svc.CreatePost(t.Context(), input)

	require.NoError(t, err)
	assert.Equal(t, expected, post)
}

// ───── CreateComment ─────

func TestCreateComment_TextTooLong(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateCommentInput{
		Text:   strings.Repeat("y", 5001),
		PostID: 1,
	}

	_, err := svc.CreateComment(t.Context(), input)

	assert.ErrorIs(t, err, ErrValidation)
}

func TestCreateComment_PostNotFound(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateCommentInput{Text: "text", PostID: 99}
	repo.EXPECT().GetPost(t.Context(), 99).Return(model.Post{}, repository.ErrNotFound)

	_, err := svc.CreateComment(t.Context(), input)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestCreateComment_RepoError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateCommentInput{Text: "text", PostID: 1}
	repo.EXPECT().GetPost(t.Context(), 1).Return(model.Post{}, assert.AnError)

	_, err := svc.CreateComment(t.Context(), input)

	assert.ErrorIs(t, err, assert.AnError)
}

func TestCreateComment_Success(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateCommentInput{Text: "text", PostID: 1}
	expected := model.Comment{ID: 1, Text: "text", PostID: 1}

	repo.EXPECT().GetPost(t.Context(), 1).Return(model.Post{ID: 1}, nil)
	repo.EXPECT().CreateComment(t.Context(), input).Return(expected, nil)

	comment, err := svc.CreateComment(t.Context(), input)

	require.NoError(t, err)
	assert.Equal(t, expected, comment)
}

// ───── GetPosts — проверка существования борда ─────

func TestGetPosts_BoardNotFound(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().GetBoard(t.Context(), 99).Return(model.Board{}, repository.ErrNotFound)

	_, err := svc.GetPosts(t.Context(), 99, false, 20, 0)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestGetPosts_BoardRepoError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().GetBoard(t.Context(), 1).Return(model.Board{}, assert.AnError)

	_, err := svc.GetPosts(t.Context(), 1, false, 20, 0)

	assert.ErrorIs(t, err, assert.AnError)
}

func TestGetPosts_NormalizesLimit(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		offset         int
		expectedLimit  int
		expectedOffset int
	}{
		{"limit zero", 0, 0, 100, 0},
		{"limit negative", -1, 0, 100, 0},
		{"limit too high", 200, 0, 100, 0},
		{"offset negative", 20, -5, 20, 0},
		{"valid", 20, 10, 20, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(t)
			svc := New(repo)

			repo.EXPECT().GetBoard(t.Context(), 1).Return(model.Board{ID: 1}, nil)
			repo.EXPECT().
				GetPosts(t.Context(), 1, false, tt.expectedLimit, tt.expectedOffset).
				Return([]model.Post{}, nil)

			_, err := svc.GetPosts(t.Context(), 1, false, tt.limit, tt.offset)

			require.NoError(t, err)
		})
	}
}

// ───── GetComments — проверка существования поста ─────

func TestGetComments_PostNotFound(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().GetPost(t.Context(), 99).Return(model.Post{}, repository.ErrNotFound)

	_, err := svc.GetComments(t.Context(), 99, false, 20, 0)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestGetComments_PostRepoError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().GetPost(t.Context(), 1).Return(model.Post{}, assert.AnError)

	_, err := svc.GetComments(t.Context(), 1, false, 20, 0)

	assert.ErrorIs(t, err, assert.AnError)
}

func TestGetComments_NormalizesLimit(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		offset         int
		expectedLimit  int
		expectedOffset int
	}{
		{"limit zero", 0, 0, 100, 0},
		{"limit negative", -1, 0, 100, 0},
		{"limit too high", 200, 0, 100, 0},
		{"offset negative", 20, -5, 20, 0},
		{"valid", 20, 10, 20, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository(t)
			svc := New(repo)

			repo.EXPECT().GetPost(t.Context(), 1).Return(model.Post{ID: 1}, nil)
			repo.EXPECT().
				GetComments(t.Context(), 1, false, tt.expectedLimit, tt.expectedOffset).
				Return([]model.Comment{}, nil)

			_, err := svc.GetComments(t.Context(), 1, false, tt.limit, tt.offset)

			require.NoError(t, err)
		})
	}
}
