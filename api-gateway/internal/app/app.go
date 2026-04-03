package app

import (
	"context"
	"fmt"

	"github.com/gfdmit/web-forum/api-gateway/config"
	v1 "github.com/gfdmit/web-forum/api-gateway/internal/handlers/http/v1"
	"github.com/gfdmit/web-forum/api-gateway/internal/httpserver"
)

func Run(conf *config.Config) error {
	ctx := context.Background()
	handler, err := v1.New(conf)
	if err != nil {
		return fmt.Errorf("error when setting up handler: %v", err)
	}

	httpserver := httpserver.New(conf.HTTPServer, handler)

	return httpserver.Run(ctx)
}
