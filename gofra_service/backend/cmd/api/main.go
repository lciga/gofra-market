// @title Gofra Market API
// @version 1.0
// @description Attack-Defense CTF сервис магазина гоферов с преднамеренными уязвимостями (NoSQL injection, SSRF, uint underflow).
// @BasePath /
// @schemes http
// @securityDefinitions.apiKey CookieAuth
// @in cookie
// @name sid
package main

import (
	"Gofra_Market/internal/app"
	"Gofra_Market/internal/config"
	"Gofra_Market/internal/db"
	"Gofra_Market/internal/docs"
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

// Точка входа программы
func main() {
	cfg := config.Load() // Загрузка конфига

	// Контекст для отмены миграции по таймауту
	migrationCtx, migrationCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer migrationCancel()

	// Запуск миграций
	if err := db.Migrate(migrationCtx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(2)
	}

	// Подключение к БД
	client, database, err := db.Connect(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db connect failed: %v\n", err)
		os.Exit(2)
	}
	defer db.Close(client)

	// Контекст для отмены сида по таймауту
	seedCtx, seedCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer seedCancel()

	// Загрузка начальных данных
	if err := db.SeedInitialData(seedCtx, database); err != nil {
		fmt.Fprintf(os.Stderr, "seed failed: %v\n", err)
		os.Exit(2)
	}

	// Инициализация коллекций
	usersColl := database.Collection("users")       // Пользователи
	gofersColl := database.Collection("gofers")     // Гоферы (лоты)
	listingsColl := database.Collection("listings") // Листинги
	sessionsColl := database.Collection("sessions") // Сессии

	// Инициализация репозиториев
	userRepo := repo.NewUserRepo(usersColl)
	goferRepo := repo.NewGoferRepo(gofersColl)
	listingRepo := repo.NewListingRepo(listingsColl)
	sessionRepo := repo.NewSessionRepo(sessionsColl)

	// Инициализация сервисов
	authSvc := service.NewAuthService(userRepo, sessionRepo, "sid")
	listingSvc := service.NewListingService(userRepo, goferRepo, listingRepo)
	marketSvc := service.NewMarketService(listingRepo, goferRepo)
	imageSvc := service.NewImageService(listingRepo)
	statsSvc := service.NewStatisticsService(sessionRepo)
	adminSvc := service.NewAdminService(userRepo, sessionRepo)
	contentSvc := service.NewContentService(userRepo, cfg.EditorEmail)

	// Инициализация хэндлеров
	authH := handlers.NewAuthHandler(authSvc, "sid")
	listingH := handlers.NewListingHandler(listingSvc)
	marketH := handlers.NewMarketHandler(marketSvc)
	imageH := handlers.NewImageHandler(imageSvc)
	statsH := handlers.NewStatisticsHandler(statsSvc)
	adminH := handlers.NewAdminHandler(adminSvc)
	contentH := handlers.NewContentHandler(contentSvc)

	engine := app.NewServer(cfg)
	engine.Use(midleware.Auth(sessionRepo))

	// Регистрация роутов
	app.RegisterRoutes(engine, app.Handlers{
		Auth:    authH,
		Market:  marketH,
		Listing: listingH,
		Image:   imageH,
		Stats:   statsH,
		Admin:   adminH,
		Content: contentH,
	})

	docs.Register(engine)

	// Запуск сервера
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
