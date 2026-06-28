package main

import (
	"log"

	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("app init: %v", err)
	}
	defer application.Close()

	if err := server.Run(application); err != nil {
		log.Printf("server: %v", err)
	}
}
