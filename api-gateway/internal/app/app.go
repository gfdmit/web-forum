package app

import (
	"context"
	"fmt"

	"github.com/gfdmit/web-forum/api-gateway/config"
	"github.com/gfdmit/web-forum/api-gateway/internal/handler"
	"github.com/gfdmit/web-forum/api-gateway/internal/httpserver"
)

func Run(conf *config.Config) error {
	ctx := context.Background()
	router, err := handler.New(conf)
	if err != nil {
		return fmt.Errorf("handler.New: %w", err)
	}

	server := httpserver.New(conf.HTTPServer, router)

	return server.Run(ctx)
}
