package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, filename string, src io.Reader, size int64, contentType string) (string, error)
	Get(ctx context.Context, filename string) (io.ReadCloser, error)
}
