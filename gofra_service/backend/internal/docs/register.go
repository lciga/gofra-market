package docs

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Gofra_Market/internal/docs/swagger"
)

// Подключает Swagger UI и служебные эндпоинты документации (только в debug-режиме).
func Register(engine *gin.Engine) {
	if gin.Mode() != gin.DebugMode {
		return
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
