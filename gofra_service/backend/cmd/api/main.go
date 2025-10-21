package main

import (
	"Gofra_Market/internal/app"
	"Gofra_Market/internal/config"
	"Gofra_Market/internal/db"
	"Gofra_Market/internal/repo"
	"Gofra_Market/internal/service"
	"Gofra_Market/internal/transport/http/handlers"
	"Gofra_Market/internal/transport/http/midleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.Load()

	if err := db.Migrate(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
	}

	client, database, err := db.Connect(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db connect failed: %v\n", err)
		os.Exit(2)
	}
	defer db.Close(client)

	if err := db.SeedInitialData(ctx, database); err != nil {
		fmt.Fprintf(os.Stderr, "seed failed: %v\n", err)
	}

	usersColl := database.Collection("users")
	gofersColl := database.Collection("gofers")
	listingsColl := database.Collection("listings")
	sessionsColl := database.Collection("sessions")

	userRepo := repo.NewUserRepo(usersColl)
	goferRepo := repo.NewGoferRepo(gofersColl)
	listingRepo := repo.NewListingRepo(listingsColl)
	sessionRepo := repo.NewSessionRepo(sessionsColl)

	authSvc := service.NewAuthService(userRepo, sessionRepo, "sid")
	listingSvc := service.NewListingService(userRepo, goferRepo, listingRepo)
	marketSvc := service.NewMarketService(listingRepo, goferRepo)
	imageSvc := service.NewImageService(listingRepo)

	authH := handlers.NewAuthHandler(authSvc, "sid")
	listingH := handlers.NewListingHandler(listingSvc)
	marketH := handlers.NewMarketHandler(marketSvc)
	imageH := handlers.NewImageHandler(imageSvc)

	engine := app.NewServer(cfg)
	engine.Use(midleware.Auth(sessionRepo))

	app.RegisterRoutes(engine, app.Handlers{
		Auth:    authH,
		Market:  marketH,
		Listing: listingH,
		Image:   imageH,
	})

	port := cfg.ServerPort
	if port == 0 {
		port = 8080
	}
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: engine}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(3)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}
