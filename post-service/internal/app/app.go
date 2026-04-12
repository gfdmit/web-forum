package app

import (
	"context"
	"fmt"

	"github.com/gfdmit/web-forum/post-service/config"
	"github.com/gfdmit/web-forum/post-service/internal/handler/rest"
	"github.com/gfdmit/web-forum/post-service/internal/httpserver"
	"github.com/gfdmit/web-forum/post-service/internal/repository/postgres"
	"github.com/gfdmit/web-forum/post-service/internal/service"
)

func Run(conf *config.Config) error {
	ctx := context.Background()

	repo, err := postgres.New(ctx, &conf.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	svc := service.New(repo)
	router := rest.NewRouter(svc)
	server := httpserver.New(conf.HTTPServer, router)

	return server.Run(ctx)
}
