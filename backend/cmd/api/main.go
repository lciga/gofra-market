package main

import (
	"Gofra_Market/internal/config"
	"Gofra_Market/internal/db"
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := config.Load()

	db.Migrate(ctx, cfg)
}
