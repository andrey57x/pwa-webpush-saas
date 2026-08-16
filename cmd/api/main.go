package main

import (
	"log"

	"github.com/andrey57x/pwa-webpush-saas/internal/app"
	"github.com/andrey57x/pwa-webpush-saas/internal/config"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("App runtime error: %v", err)
	}
}
