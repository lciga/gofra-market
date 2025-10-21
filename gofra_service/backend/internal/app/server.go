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

	gin.SetMode(mode)
	logger.Infof("Gin server start in mode: %d", port)

	eng := gin.New()

	eng.Use(gin.Logger())
	eng.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if len(cfg.AllowedOrigins) > 0 && cfg.AllowedOrigins[0] == "*" {
				return true
			}
			for _, allowed := range cfg.AllowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	eng.Use(cors.New(corsConfig))

	return eng
}
