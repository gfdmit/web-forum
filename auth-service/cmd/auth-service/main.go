package main

import (
	"log"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/app"
)

func main() {
	conf, err := config.New(".env")
	if err != nil {
		log.Fatalf("[SETUP ERROR] error when reading config: %v", err)
	}

	err = app.Run(conf)
	if err != nil {
		log.Fatalf("[APPLICATION ERROR] error: %v", err)
	}

	log.Println("[SHUTDOWN] service shut down gracefully")
}
