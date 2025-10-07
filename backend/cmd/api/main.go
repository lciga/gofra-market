package main

import (
	"Gofra_Market/internal/app"
	"Gofra_Market/internal/config"
)

func main() {
	cfg := config.Load()

	eng := app.NewServer(cfg)
	eng.Run()
}
