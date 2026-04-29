package rest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gfdmit/web-forum/media-service/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	store      storage.Storage
	publicHost string
}

func New(store storage.Storage, publicHost string) *Handler {
	return &Handler{store: store, publicHost: publicHost}
}

func (h *Handler) PostMedia(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFileSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, errorResponse(fmt.Errorf("файл превышает 10 МБ")))
			return
		}
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("поле 'file' обязательно")))
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		handleError(c, fmt.Errorf("file.Read: %w", err))
		return
	}
	contentType := http.DetectContentType(buf[:n])

	ext, ok := allowedTypes[contentType]
	if !ok {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("недопустимый тип файла: %s", contentType)))
		return
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		handleError(c, fmt.Errorf("file.Seek: %w", err))
		return
	}

	filename := uuid.NewString() + ext

	if _, err = h.store.Upload(c.Request.Context(), filename, file, header.Size, contentType); err != nil {
		handleError(c, fmt.Errorf("store.Upload: %w", err))
		return
	}

	url := fmt.Sprintf("%s/api/v1/media/%s", h.publicHost, filename)
	c.JSON(http.StatusCreated, gin.H{"url": url})
}

func (h *Handler) GetMedia(c *gin.Context) {
	filename := c.Param("filename")

	if strings.Contains(filename, "/") || strings.Contains(filename, "..") {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("недопустимое имя файла")))
		return
	}

	obj, err := h.store.Get(c.Request.Context(), filename)
	if err != nil {
		handleError(c, fmt.Errorf("store.Get: %w", err))
		return
	}
	defer obj.Close()

	contentType := extensionToContentType(filepath.Ext(filename))
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=31536000")

	if _, err = io.Copy(c.Writer, obj); err != nil {
		log.Printf("[ERROR] GetMedia io.Copy: %v", err)
	}
}
