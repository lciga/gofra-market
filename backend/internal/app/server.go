package app

import (
	"Gofra_Market/internal/config"
	"Gofra_Market/internal/logger"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config) *gin.Engine {
	mode := cfg.GinMode
	port := cfg.ServerPort

	// Устанавливаем режим Gin
	gin.SetMode(mode)
	logger.Infof("Gin server start in mode: %d", port)

	// Создаем новый сервер Gin
	eng := gin.New()

	// Добавляем middleware для логирования и восстановления после паники
	eng.Use(gin.Logger())
	eng.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	eng.Use(cors.New(corsConfig))

	return eng
}
