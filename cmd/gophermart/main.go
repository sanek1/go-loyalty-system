package main

import (
	"context"
	"go-loyalty-system/config"
	"go-loyalty-system/internal/app"
	"log"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Setup error: %s", err)
	}
	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("Run error: %s", err)
	}
}
