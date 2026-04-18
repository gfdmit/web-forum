package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gfdmit/web-forum/post-service/internal/mocks"
	"github.com/gfdmit/web-forum/post-service/internal/model"
	"github.com/gfdmit/web-forum/post-service/internal/repository"
	"github.com/gfdmit/web-forum/post-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

// ───── parseID — покрываем один раз для всех эндпоинтов с :id ─────

func TestParseID_Invalid(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v2/boards/abc"},
		{http.MethodDelete, "/api/v2/boards/abc"},
		{http.MethodPost, "/api/v2/boards/abc/restore"},
		{http.MethodGet, "/api/v2/posts/abc"},
		{http.MethodDelete, "/api/v2/posts/abc"},
		{http.MethodGet, "/api/v2/comments/abc"},
		{http.MethodDelete, "/api/v2/comments/abc"},
		{http.MethodGet, "/api/v2/profiles/abc"},
	}

	for _, e := range endpoints {
		t.Run(e.method+" "+e.path, func(t *testing.T) {
			svc := mocks.NewMockService(t)
			router := NewRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(e.method, e.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

// ───── parseUserID — покрываем один раз через CreatePost ─────

func TestParseUserID_Missing(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/posts", jsonBody(t, model.CreatePostInput{}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParseUserID_Invalid(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/posts", jsonBody(t, model.CreatePostInput{}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "not-a-number")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParseUserID_Zero(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/posts", jsonBody(t, model.CreatePostInput{}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "0")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ───── GetBoards ─────

func TestGetBoards_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetBoards(t.Context(), false).Return([]model.Board{{ID: 1, Name: "General"}}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/boards", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBoards_ServiceError(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetBoards(t.Context(), false).Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/boards", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ───── GetBoard ─────

func TestGetBoard_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetBoard(t.Context(), 99).Return(model.Board{}, repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/boards/99", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetBoard_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetBoard(t.Context(), 1).Return(model.Board{ID: 1, Name: "General"}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/boards/1", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ───── CreateBoard ─────

func TestCreateBoard_InvalidJSON(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/boards", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBoard_ValidationError(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	input := model.CreateBoardInput{Name: ""}
	svc.EXPECT().CreateBoard(t.Context(), input).Return(model.Board{}, service.ErrValidation)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/boards", jsonBody(t, input))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBoard_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	input := model.CreateBoardInput{Name: "General"}
	svc.EXPECT().CreateBoard(t.Context(), input).Return(model.Board{ID: 1, Name: "General"}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/boards", jsonBody(t, input))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// ───── DeleteBoard ─────

func TestDeleteBoard_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().DeleteBoard(t.Context(), 99).Return(repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/boards/99", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBoard_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().DeleteBoard(t.Context(), 1).Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/boards/1", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ───── RestoreBoard ─────

func TestRestoreBoard_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().RestoreBoard(t.Context(), 99).Return(repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/boards/99/restore", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRestoreBoard_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().RestoreBoard(t.Context(), 1).Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/boards/1/restore", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ───── GetPost ─────

func TestGetPost_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetPost(t.Context(), 99).Return(model.Post{}, repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/posts/99", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPost_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetPost(t.Context(), 1).Return(model.Post{ID: 1, BoardID: 1}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/posts/1", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ───── GetPosts ─────

func TestGetPosts_QueryParams(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		includeDeleted bool
		limit          int
		offset         int
	}{
		{"defaults", "", false, 20, 0},
		{"includeDeleted", "?includeDeleted=true", true, 20, 0},
		{"custom limit offset", "?limit=10&offset=5", false, 10, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockService(t)
			router := NewRouter(svc)

			svc.EXPECT().
				GetPosts(t.Context(), 1, tt.includeDeleted, tt.limit, tt.offset).
				Return([]model.Post{}, nil)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v2/boards/1/posts"+tt.query, nil)
			req = req.WithContext(t.Context())
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestGetPosts_ServiceError(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetPosts(t.Context(), 1, false, 20, 0).Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/boards/1/posts", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ───── CreatePost ─────

func TestCreatePost_InvalidJSON(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/posts", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePost_BoardNotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	userID := 1
	input := model.CreatePostInput{BoardID: 99, Text: "text", UserID: &userID}
	svc.EXPECT().CreatePost(t.Context(), input).Return(model.Post{}, repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/posts", jsonBody(t, input))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreatePost_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	userID := 1
	input := model.CreatePostInput{BoardID: 1, Text: "text", UserID: &userID}
	svc.EXPECT().CreatePost(t.Context(), input).Return(model.Post{ID: 1, BoardID: 1, Text: "text"}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/posts", jsonBody(t, input))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// ───── DeletePost ─────

func TestDeletePost_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().DeletePost(t.Context(), 99).Return(repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/posts/99", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeletePost_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().DeletePost(t.Context(), 1).Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/posts/1", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ───── GetComment ─────

func TestGetComment_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetComment(t.Context(), 99).Return(model.Comment{}, repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/comments/99", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetComment_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetComment(t.Context(), 1).Return(model.Comment{ID: 1, PostID: 1}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/comments/1", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ───── GetComments ─────

func TestGetComments_QueryParams(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		includeDeleted bool
		limit          int
		offset         int
	}{
		{"defaults", "", false, 20, 0},
		{"includeDeleted", "?includeDeleted=true", true, 20, 0},
		{"custom limit offset", "?limit=10&offset=5", false, 10, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockService(t)
			router := NewRouter(svc)

			svc.EXPECT().
				GetComments(t.Context(), 1, tt.includeDeleted, tt.limit, tt.offset).
				Return([]model.Comment{}, nil)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v2/posts/1/comments"+tt.query, nil)
			req = req.WithContext(t.Context())
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestGetComments_ServiceError(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetComments(t.Context(), 1, false, 20, 0).Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/posts/1/comments", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ───── CreateComment ─────

func TestCreateComment_InvalidJSON(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/comments", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateComment_PostNotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	userID := 1
	input := model.CreateCommentInput{PostID: 99, Text: "text", UserID: &userID}
	svc.EXPECT().CreateComment(t.Context(), input).Return(model.Comment{}, repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/comments", jsonBody(t, input))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateComment_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	userID := 1
	input := model.CreateCommentInput{PostID: 1, Text: "text", UserID: &userID}
	svc.EXPECT().CreateComment(t.Context(), input).Return(model.Comment{ID: 1, PostID: 1, Text: "text"}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/comments", jsonBody(t, input))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// ───── DeleteComment ─────

func TestDeleteComment_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().DeleteComment(t.Context(), 99).Return(repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/comments/99", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteComment_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().DeleteComment(t.Context(), 1).Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/comments/1", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ───── GetProfiles ─────

func TestGetProfiles_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetProfiles(t.Context(), false).Return([]model.Profile{}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/profiles", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProfiles_ServiceError(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetProfiles(t.Context(), false).Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/profiles", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ───── GetProfile ─────

func TestGetProfile_NotFound(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetProfile(t.Context(), 99).Return(model.Profile{}, repository.ErrNotFound)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/profiles/99", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetProfile_Success(t *testing.T) {
	svc := mocks.NewMockService(t)
	router := NewRouter(svc)

	svc.EXPECT().GetProfile(t.Context(), 1).Return(model.Profile{ID: 1, UserID: 1}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/profiles/1", nil)
	req = req.WithContext(t.Context())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
