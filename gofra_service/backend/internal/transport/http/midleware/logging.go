package midleware

import "github.com/gin-gonic/gin"

// Логирование попыток аутентификации
func Logging() gin.HandlerFunc {
	return gin.Logger()
}
