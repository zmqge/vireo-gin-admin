package middleware

import (
	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")                                    // 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")     // 允许的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization") // 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")                            // 允许携带凭证

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
