package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/config"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	// 允许的域名列表
	allowedOrigins := config.App.AllowedOrigins
	allowedOriginsMap := make(map[string]bool)

	for _, origin := range allowedOrigins {
		allowedOriginsMap[origin] = true
	}

	log.Printf("CORS中间件初始化: 允许的域名=%v", allowedOrigins)

	return func(c *gin.Context) {
		// 获取请求的Origin头
		origin := c.Request.Header.Get("Origin")

		// 允许的请求头列表（根据实际需求调整）
		allowedHeaders := "Origin, Content-Type, Accept, Authorization, X-Requested-With"

		// 允许的HTTP方法
		allowedMethods := "GET, POST, PUT, DELETE, OPTIONS"

		// 检查Origin是否在允许列表中
		if origin != "" && allowedOriginsMap[origin] {
			// 只对允许的域名设置响应头
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")

			log.Printf("CORS允许请求: Origin=%s, Path=%s, Method=%s",
				origin, c.Request.URL.Path, c.Request.Method)
		} else {
			log.Printf("CORS拒绝请求: Origin=%s, Path=%s, Method=%s",
				origin, c.Request.URL.Path, c.Request.Method)

			// 非预检请求且Origin不在允许列表，返回403
			if c.Request.Method != "OPTIONS" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code": 403,
					"msg":  "跨域请求被拒绝",
				})
				return
			}
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			// 记录预检请求信息
			log.Printf("处理预检请求: Origin=%s, Path=%s, Access-Control-Request-Method=%s",
				origin, c.Request.URL.Path, c.Request.Header.Get("Access-Control-Request-Method"))

			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
