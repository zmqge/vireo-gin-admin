package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.Next()
	}
}
