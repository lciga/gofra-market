package midleware

import "github.com/gin-gonic/gin"

func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
	c.Abort()
}
