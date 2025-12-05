package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

	root := projectRoot()
	handler, err := NewPackageDocsHandler(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "package docs disabled: %v\n", err)
		return
	}

	engine.GET("/debug/pkgdocs", handler)
}

func projectRoot() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	}
	return "."
}
