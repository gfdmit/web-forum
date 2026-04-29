package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gfdmit/web-forum/media-service/config"
	"github.com/gfdmit/web-forum/media-service/internal/storage"
)

type minioStorage struct {
	cli    *minio.Client
	bucket string
}

func New(conf config.MinIO) (storage.Storage, error) {
	client, err := minio.New(fmt.Sprintf("%s:%s", conf.Host, conf.Port), &minio.Options{
		Creds:  credentials.NewStaticV4(conf.User, conf.Pass, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("minio.New: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, conf.Bucket)
	if err != nil {
		return nil, fmt.Errorf("client.BucketExists: %w", err)
	}

	if !exists {
		if err = client.MakeBucket(ctx, conf.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("client.MakeBucket: %w", err)
		}
	}

	return &minioStorage{cli: client, bucket: conf.Bucket}, nil
}

func (ms *minioStorage) Upload(ctx context.Context, filename string, src io.Reader, size int64, contentType string) (string, error) {
	_, err := ms.cli.PutObject(ctx, ms.bucket, filename, src, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("cli.PutObject: %w", err)
	}

	return filename, nil
}

func (ms *minioStorage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	obj, err := ms.cli.GetObject(ctx, ms.bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("cli.GetObject: %w", err)
	}

	if _, err = obj.Stat(); err != nil {
		obj.Close()
		return nil, fmt.Errorf("obj.Stat: %w", err)
	}

	return obj, nil
}
