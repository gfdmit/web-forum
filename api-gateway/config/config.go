package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	PostService
	//MediaService
	HTTPServer
}

type PostService struct {
	Host string `env:"POSTSERVICE_HOST" env-default:"localhost"`
	Port string `env:"POSTSERVICE_PORT" env-default:"8080"`
}

type HTTPServer struct {
	BindAddress     string        `env:"BIND_ADDRESS" env-default:"localhost"`
	BindPort        string        `env:"BIND_PORT" env-default:"3030"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"5s"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" env-default:"5s"`
}

func New(env string) (*Config, error) {
	conf := &Config{}

	if err := godotenv.Overload(env); err != nil {
		return nil, fmt.Errorf("godotenv.Overload: %v", err)
	}

	if err := cleanenv.ReadEnv(conf); err != nil {
		return nil, fmt.Errorf("cleanenv.Readenv: %v", err)
	}

	return conf, nil
}
