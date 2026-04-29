package app

import (
	"context"
	"fmt"

	"github.com/gfdmit/web-forum/media-service/config"
	"github.com/gfdmit/web-forum/media-service/internal/handler/rest"
	"github.com/gfdmit/web-forum/media-service/internal/httpserver"
	"github.com/gfdmit/web-forum/media-service/internal/storage/minio"
)

func Run(conf config.Config) error {
	ctx := context.Background()

	store, err := minio.New(conf.MinIO)
	if err != nil {
		return fmt.Errorf("error when setting up repository: %v", err)
	}

	handler := rest.NewRouter(store, conf.MinIO.PublicHost)

	server := httpserver.New(conf.HTTPServer, handler)

	return server.Run(ctx)
}
