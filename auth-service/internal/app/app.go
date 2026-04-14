package app

import (
	"context"
	"fmt"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/handler/rest"
	"github.com/gfdmit/web-forum/auth-service/internal/httpserver"
	"github.com/gfdmit/web-forum/auth-service/internal/kfu"
	"github.com/gfdmit/web-forum/auth-service/internal/repository/postgres"
	"github.com/gfdmit/web-forum/auth-service/internal/service"
)

func Run(conf *config.Config) error {
	ctx := context.Background()

	repo, err := postgres.New(ctx, conf.Postgres)
	if err != nil {
		return fmt.Errorf("error when setting up repository: %w", err)
	}

	client := kfu.NewClient()

	svc := service.New(&conf.JWT, repo, client)

	handler := rest.NewRouter(svc, conf.JWT.TTL)

	server := httpserver.New(conf.HTTPServer, handler)

	return server.Run(ctx)
}
