//go:build integration

package postgres_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/gfdmit/web-forum/post-service/config"
	"github.com/gfdmit/web-forum/post-service/internal/model"
	"github.com/gfdmit/web-forum/post-service/internal/repository"
	postgres "github.com/gfdmit/web-forum/post-service/internal/repository/postgres"
)

func setupTestDB(t *testing.T) repository.Repository {
	t.Helper()
	ctx := t.Context()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:18-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "forum_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(context.Background()) })

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "../../../migrations")

	conf := &config.Postgres{
		User:       "test",
		Pass:       "test",
		Host:       host,
		Port:       port.Port(),
		DB:         "forum_test",
		Migrations: migrationsPath,
		Pool: config.Pool{
			MaxConns:          5,
			MinConns:          1,
			HealthCheckPeriod: 30 * time.Second,
		},
	}

	repo, err := postgres.New(ctx, conf)
	require.NoError(t, err)

	return repo
}

func TestRepo_Board(t *testing.T) {
	repo := setupTestDB(t)
	ctx := t.Context()

	t.Run("CreateBoard_Success", func(t *testing.T) {
		desc := "a fresh board"
		board, err := repo.CreateBoard(ctx, model.CreateBoardInput{
			Name:        "General",
			Description: &desc,
		})
		require.NoError(t, err)
		assert.Greater(t, board.ID, 0)
		assert.Equal(t, "General", board.Name)
		assert.Equal(t, &desc, board.Description)
		assert.Nil(t, board.DeletedAt)
	})

	t.Run("GetBoard_Active", func(t *testing.T) {
		board := createBoard(t, repo, "Active Board")
		got, err := repo.GetBoard(ctx, board.ID)
		require.NoError(t, err)
		assert.Equal(t, board.ID, got.ID)
		assert.Equal(t, "Active Board", got.Name)
	})

	t.Run("GetBoard_NonExistent", func(t *testing.T) {
		_, err := repo.GetBoard(ctx, 999999)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetBoard_Deleted_ReturnsNotFound", func(t *testing.T) {
		board := createBoard(t, repo, "To Delete")
		require.NoError(t, repo.DeleteBoard(ctx, board.ID))

		_, err := repo.GetBoard(ctx, board.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetBoards_ExcludesDeleted", func(t *testing.T) {
		active := createBoard(t, repo, "Visible Board")
		deleted := createBoard(t, repo, "Hidden Board")
		require.NoError(t, repo.DeleteBoard(ctx, deleted.ID))

		boards, err := repo.GetBoards(ctx, false)
		require.NoError(t, err)

		ids := make([]int, 0, len(boards))
		for _, b := range boards {
			ids = append(ids, b.ID)
		}
		assert.Contains(t, ids, active.ID)
		assert.NotContains(t, ids, deleted.ID)
	})

	t.Run("GetBoards_IncludesDeleted", func(t *testing.T) {
		deleted := createBoard(t, repo, "Ghost Board")
		require.NoError(t, repo.DeleteBoard(ctx, deleted.ID))

		boards, err := repo.GetBoards(ctx, true)
		require.NoError(t, err)

		ids := make([]int, 0, len(boards))
		for _, b := range boards {
			ids = append(ids, b.ID)
		}
		assert.Contains(t, ids, deleted.ID)
	})

	t.Run("DeleteBoard_Success", func(t *testing.T) {
		board := createBoard(t, repo, "Board To Delete")
		err := repo.DeleteBoard(ctx, board.ID)
		require.NoError(t, err)

		_, err = repo.GetBoard(ctx, board.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("DeleteBoard_NonExistent", func(t *testing.T) {
		err := repo.DeleteBoard(ctx, 999998)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("RestoreBoard_Success", func(t *testing.T) {
		board := createBoard(t, repo, "Board To Restore")
		require.NoError(t, repo.DeleteBoard(ctx, board.ID))

		err := repo.RestoreBoard(ctx, board.ID)
		require.NoError(t, err)

		got, err := repo.GetBoard(ctx, board.ID)
		require.NoError(t, err)
		assert.Equal(t, board.ID, got.ID)
		assert.Nil(t, got.DeletedAt)
	})

	t.Run("RestoreBoard_NonExistent", func(t *testing.T) {
		err := repo.RestoreBoard(ctx, 999997)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestRepo_Post(t *testing.T) {
	repo := setupTestDB(t)
	ctx := t.Context()

	t.Run("CreatePost_Success_NoUser", func(t *testing.T) {
		board := createBoard(t, repo, "Post Board")
		post, err := repo.CreatePost(ctx, model.CreatePostInput{
			UserID:  nil,
			BoardID: board.ID,
			Title:   "Hello World",
			Text:    "First post",
		})
		require.NoError(t, err)
		assert.Greater(t, post.ID, 0)
		assert.Equal(t, "Hello World", post.Title)
		assert.Nil(t, post.UserID)
		assert.Nil(t, post.Author)
	})

	t.Run("GetPost_Active", func(t *testing.T) {
		board := createBoard(t, repo, "GetPost Board")
		post := createPost(t, repo, board.ID, nil)

		got, err := repo.GetPost(ctx, post.ID)
		require.NoError(t, err)
		assert.Equal(t, post.ID, got.ID)
		assert.Equal(t, board.ID, got.BoardID)
	})

	t.Run("GetPost_NonExistent", func(t *testing.T) {
		_, err := repo.GetPost(ctx, 999999)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetPost_Deleted_ReturnsNotFound", func(t *testing.T) {
		board := createBoard(t, repo, "GetPost Deleted Board")
		post := createPost(t, repo, board.ID, nil)
		require.NoError(t, repo.DeletePost(ctx, post.ID))

		_, err := repo.GetPost(ctx, post.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetPosts_ExcludesDeleted", func(t *testing.T) {
		board := createBoard(t, repo, "GetPosts Board")
		active := createPost(t, repo, board.ID, nil)
		deleted := createPost(t, repo, board.ID, nil)
		require.NoError(t, repo.DeletePost(ctx, deleted.ID))

		posts, err := repo.GetPosts(ctx, board.ID, false, 100, 0)
		require.NoError(t, err)

		ids := postIDs(posts)
		assert.Contains(t, ids, active.ID)
		assert.NotContains(t, ids, deleted.ID)
	})

	t.Run("GetPosts_IncludesDeleted", func(t *testing.T) {
		board := createBoard(t, repo, "GetPosts IncludeDel Board")
		deleted := createPost(t, repo, board.ID, nil)
		require.NoError(t, repo.DeletePost(ctx, deleted.ID))

		posts, err := repo.GetPosts(ctx, board.ID, true, 100, 0)
		require.NoError(t, err)

		ids := postIDs(posts)
		assert.Contains(t, ids, deleted.ID)
	})

	t.Run("GetPosts_Pagination", func(t *testing.T) {
		board := createBoard(t, repo, "Pagination Board")
		for i := 0; i < 5; i++ {
			createPost(t, repo, board.ID, nil)
		}

		page1, err := repo.GetPosts(ctx, board.ID, false, 3, 0)
		require.NoError(t, err)
		assert.Len(t, page1, 3)

		page2, err := repo.GetPosts(ctx, board.ID, false, 3, 3)
		require.NoError(t, err)
		assert.Len(t, page2, 2)
	})

	t.Run("DeletePost_Success", func(t *testing.T) {
		board := createBoard(t, repo, "DeletePost Board")
		post := createPost(t, repo, board.ID, nil)

		require.NoError(t, repo.DeletePost(ctx, post.ID))
		_, err := repo.GetPost(ctx, post.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("DeletePost_NonExistent", func(t *testing.T) {
		err := repo.DeletePost(ctx, 999998)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("Post_Author_NilWhenNoUserID", func(t *testing.T) {
		board := createBoard(t, repo, "Author Nil Board")
		post := createPost(t, repo, board.ID, nil)

		got, err := repo.GetPost(ctx, post.ID)
		require.NoError(t, err)
		assert.Nil(t, got.Author)
	})
}

func TestRepo_Comment(t *testing.T) {
	repo := setupTestDB(t)
	ctx := t.Context()

	t.Run("CreateComment_Success", func(t *testing.T) {
		board := createBoard(t, repo, "Comment Board")
		post := createPost(t, repo, board.ID, nil)

		comment, err := repo.CreateComment(ctx, model.CreateCommentInput{
			UserID: nil,
			PostID: post.ID,
			Text:   "First comment",
		})
		require.NoError(t, err)
		assert.Greater(t, comment.ID, 0)
		assert.Equal(t, post.ID, comment.PostID)
		assert.Equal(t, "First comment", comment.Text)
		assert.Nil(t, comment.Author)
	})

	t.Run("GetComment_Active", func(t *testing.T) {
		board := createBoard(t, repo, "GetComment Board")
		post := createPost(t, repo, board.ID, nil)
		comment := createComment(t, repo, post.ID, nil)

		got, err := repo.GetComment(ctx, comment.ID)
		require.NoError(t, err)
		assert.Equal(t, comment.ID, got.ID)
		assert.Equal(t, post.ID, got.PostID)
	})

	t.Run("GetComment_NonExistent", func(t *testing.T) {
		_, err := repo.GetComment(ctx, 999999)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetComment_Deleted_ReturnsNotFound", func(t *testing.T) {
		board := createBoard(t, repo, "GetComment Deleted Board")
		post := createPost(t, repo, board.ID, nil)
		comment := createComment(t, repo, post.ID, nil)
		require.NoError(t, repo.DeleteComment(ctx, comment.ID))

		_, err := repo.GetComment(ctx, comment.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetComments_ExcludesDeleted", func(t *testing.T) {
		board := createBoard(t, repo, "GetComments Board")
		post := createPost(t, repo, board.ID, nil)
		active := createComment(t, repo, post.ID, nil)
		deleted := createComment(t, repo, post.ID, nil)
		require.NoError(t, repo.DeleteComment(ctx, deleted.ID))

		comments, err := repo.GetComments(ctx, post.ID, false, 100, 0)
		require.NoError(t, err)

		ids := commentIDs(comments)
		assert.Contains(t, ids, active.ID)
		assert.NotContains(t, ids, deleted.ID)
	})

	t.Run("GetComments_IncludesDeleted", func(t *testing.T) {
		board := createBoard(t, repo, "GetComments IncludeDel Board")
		post := createPost(t, repo, board.ID, nil)
		deleted := createComment(t, repo, post.ID, nil)
		require.NoError(t, repo.DeleteComment(ctx, deleted.ID))

		comments, err := repo.GetComments(ctx, post.ID, true, 100, 0)
		require.NoError(t, err)

		ids := commentIDs(comments)
		assert.Contains(t, ids, deleted.ID)
	})

	t.Run("GetComments_Pagination", func(t *testing.T) {
		board := createBoard(t, repo, "Comment Pagination Board")
		post := createPost(t, repo, board.ID, nil)
		for i := 0; i < 5; i++ {
			createComment(t, repo, post.ID, nil)
		}

		page1, err := repo.GetComments(ctx, post.ID, false, 3, 0)
		require.NoError(t, err)
		assert.Len(t, page1, 3)

		page2, err := repo.GetComments(ctx, post.ID, false, 3, 3)
		require.NoError(t, err)
		assert.Len(t, page2, 2)
	})

	t.Run("DeleteComment_Success", func(t *testing.T) {
		board := createBoard(t, repo, "DeleteComment Board")
		post := createPost(t, repo, board.ID, nil)
		comment := createComment(t, repo, post.ID, nil)

		require.NoError(t, repo.DeleteComment(ctx, comment.ID))
		_, err := repo.GetComment(ctx, comment.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("DeleteComment_NonExistent", func(t *testing.T) {
		err := repo.DeleteComment(ctx, 999998)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("Comment_Author_NilWhenNoUserID", func(t *testing.T) {
		board := createBoard(t, repo, "Comment Author Nil Board")
		post := createPost(t, repo, board.ID, nil)
		comment := createComment(t, repo, post.ID, nil)

		got, err := repo.GetComment(ctx, comment.ID)
		require.NoError(t, err)
		assert.Nil(t, got.Author)
	})
}

func TestRepo_Profile(t *testing.T) {
	repo := setupTestDB(t)
	ctx := t.Context()

	t.Run("GetProfile_NonExistent", func(t *testing.T) {
		_, err := repo.GetProfile(ctx, 999999)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetProfiles_EmptyWhenNone", func(t *testing.T) {
		profiles, err := repo.GetProfiles(ctx, false)
		require.NoError(t, err)
		assert.NotNil(t, profiles)
	})
}

func postIDs(posts []model.Post) []int {
	ids := make([]int, 0, len(posts))
	for _, p := range posts {
		ids = append(ids, p.ID)
	}
	return ids
}

func commentIDs(comments []model.Comment) []int {
	ids := make([]int, 0, len(comments))
	for _, c := range comments {
		ids = append(ids, c.ID)
	}
	return ids
}

func createBoard(t *testing.T, repo repository.Repository, name string) model.Board {
	t.Helper()
	desc := "test description"
	board, err := repo.CreateBoard(t.Context(), model.CreateBoardInput{
		Name:        name,
		Description: &desc,
	})
	require.NoError(t, err)
	return board
}

func createPost(t *testing.T, repo repository.Repository, boardID int, userID *int) model.Post {
	t.Helper()
	post, err := repo.CreatePost(t.Context(), model.CreatePostInput{
		UserID:  userID,
		BoardID: boardID,
		Title:   "Test Post",
		Text:    "Test post body",
	})
	require.NoError(t, err)
	return post
}

func createComment(t *testing.T, repo repository.Repository, postID int, userID *int) model.Comment {
	t.Helper()
	comment, err := repo.CreateComment(t.Context(), model.CreateCommentInput{
		UserID: userID,
		PostID: postID,
		Text:   "Test comment",
	})
	require.NoError(t, err)
	return comment
}
