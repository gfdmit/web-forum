package service

import (
	"context"
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

	_, err := svc.CreateBoard(context.Background(), model.CreateBoardInput{Name: ""})

	assert.ErrorIs(t, err, ErrValidation)
}

func TestCreateBoard_Success(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateBoardInput{Name: "General"}
	expected := model.Board{ID: 1, Name: "General"}

	repo.EXPECT().CreateBoard(context.Background(), input).Return(expected, nil).Once()

	board, err := svc.CreateBoard(context.Background(), input)

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

	_, err := svc.CreatePost(context.Background(), input)

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

	_, err := svc.CreatePost(context.Background(), input)

	assert.ErrorIs(t, err, ErrValidation)
}

func TestCreatePost_BoardNotFound(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{Text: "text", BoardID: 99}
	repo.EXPECT().GetBoard(context.Background(), 99).Return(model.Board{}, repository.ErrNotFound).Once()

	_, err := svc.CreatePost(context.Background(), input)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestCreatePost_RepoError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{Text: "text", BoardID: 1}
	repo.EXPECT().GetBoard(context.Background(), 1).Return(model.Board{}, assert.AnError).Once()

	_, err := svc.CreatePost(context.Background(), input)

	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestCreatePost_Success(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreatePostInput{Text: "text", BoardID: 1}
	expected := model.Post{ID: 1, Text: "text", BoardID: 1}

	repo.EXPECT().GetBoard(context.Background(), 1).Return(model.Board{ID: 1}, nil).Once()
	repo.EXPECT().CreatePost(context.Background(), input).Return(expected, nil).Once()

	post, err := svc.CreatePost(context.Background(), input)

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

	_, err := svc.CreateComment(context.Background(), input)

	assert.ErrorIs(t, err, ErrValidation)
}

func TestCreateComment_PostNotFound(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateCommentInput{Text: "text", PostID: 99}
	repo.EXPECT().GetPost(context.Background(), 99).Return(model.Post{}, repository.ErrNotFound).Once()

	_, err := svc.CreateComment(context.Background(), input)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestCreateComment_Success(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateCommentInput{Text: "text", PostID: 1}
	expected := model.Comment{ID: 1, Text: "text", PostID: 1}

	repo.EXPECT().GetPost(context.Background(), 1).Return(model.Post{ID: 1}, nil).Once()
	repo.EXPECT().CreateComment(context.Background(), input).Return(expected, nil).Once()

	comment, err := svc.CreateComment(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, expected, comment)
}

// ───── GetPosts — нормализация лимитов ─────

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

			repo.EXPECT().
				GetPosts(context.Background(), 1, false, tt.expectedLimit, tt.expectedOffset).
				Return([]model.Post{}, nil).
				Once()

			_, err := svc.GetPosts(context.Background(), 1, false, tt.limit, tt.offset)

			require.NoError(t, err)
		})
	}
}

// ───── GetComments — нормализация лимитов ─────

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

			repo.EXPECT().
				GetComments(context.Background(), 1, false, tt.expectedLimit, tt.expectedOffset).
				Return([]model.Comment{}, nil).
				Once()

			_, err := svc.GetComments(context.Background(), 1, false, tt.limit, tt.offset)

			require.NoError(t, err)
		})
	}
}

// ───── GetBoard / GetBoards ─────

func TestGetBoard_PassesThroughResult(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	want := model.Board{ID: 1, Name: "General"}
	repo.EXPECT().GetBoard(context.Background(), 1).Return(want, nil).Once()

	got, err := svc.GetBoard(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetBoard_PassesThroughError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().GetBoard(context.Background(), 1).Return(model.Board{}, repository.ErrNotFound).Once()

	_, err := svc.GetBoard(context.Background(), 1)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestGetBoards_PassesThroughResult(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	want := []model.Board{{ID: 1}, {ID: 2}}
	repo.EXPECT().GetBoards(context.Background(), false).Return(want, nil).Once()

	got, err := svc.GetBoards(context.Background(), false)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// ───── DeleteBoard / RestoreBoard ─────

func TestDeleteBoard_PassesThroughError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().DeleteBoard(context.Background(), 1).Return(repository.ErrNotFound).Once()

	err := svc.DeleteBoard(context.Background(), 1)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestRestoreBoard_PassesThroughError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().RestoreBoard(context.Background(), 1).Return(repository.ErrNotFound).Once()

	err := svc.RestoreBoard(context.Background(), 1)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

// ───── GetPost / DeletePost ─────

func TestGetPost_PassesThroughResult(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	want := model.Post{ID: 5}
	repo.EXPECT().GetPost(context.Background(), 5).Return(want, nil).Once()

	got, err := svc.GetPost(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDeletePost_PassesThroughError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().DeletePost(context.Background(), 1).Return(repository.ErrNotFound).Once()

	err := svc.DeletePost(context.Background(), 1)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

// ───── GetComment / DeleteComment ─────

func TestGetComment_PassesThroughResult(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	want := model.Comment{ID: 3}
	repo.EXPECT().GetComment(context.Background(), 3).Return(want, nil).Once()

	got, err := svc.GetComment(context.Background(), 3)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDeleteComment_PassesThroughError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	repo.EXPECT().DeleteComment(context.Background(), 1).Return(repository.ErrNotFound).Once()

	err := svc.DeleteComment(context.Background(), 1)

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

// ───── GetProfile / GetProfiles ─────

func TestGetProfile_PassesThroughResult(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	want := model.Profile{UserID: 7}
	repo.EXPECT().GetProfile(context.Background(), 7).Return(want, nil).Once()

	got, err := svc.GetProfile(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetProfiles_PassesThroughResult(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	want := []model.Profile{{UserID: 1}, {UserID: 2}}
	repo.EXPECT().GetProfiles(context.Background(), false).Return(want, nil).Once()

	got, err := svc.GetProfiles(context.Background(), false)

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// ───── wrapped errors ─────

func TestCreateComment_RepoError(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	svc := New(repo)

	input := model.CreateCommentInput{Text: "text", PostID: 1}
	repo.EXPECT().GetPost(context.Background(), 1).Return(model.Post{}, assert.AnError).Once()

	_, err := svc.CreateComment(context.Background(), input)

	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
