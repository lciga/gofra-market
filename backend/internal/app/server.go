package app

import (
	"Gofra_Market/internal/config"
	"Gofra_Market/internal/logger"

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

	eng.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method != "POST" || c.Request.Method != "GET" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	return eng
}
