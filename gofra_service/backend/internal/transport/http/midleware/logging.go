package midleware

import "github.com/gin-gonic/gin"

// Логирует метод/путь/код/время.
func Logging() gin.HandlerFunc {
	return gin.Logger()
}
