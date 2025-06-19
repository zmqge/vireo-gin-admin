package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

func DemoMode() gin.HandlerFunc {
	// 允许的路径列表（支持前缀匹配和通配符）
	allowedPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/refresh_token",
		"/api/v1/auth/logout",
		"/api/v1/auth/captcha",
		// 添加更多允许的路径...
	}

	return func(c *gin.Context) {
		// 如果不是演示模式，直接放行
		if !config.App.DemoMode {
			c.Next()
			return
		}

		requestPath := c.Request.URL.Path

		// 检查请求路径是否在允许的列表中（支持前缀匹配）
		for _, path := range allowedPaths {
			// 精确匹配
			if requestPath == path {
				c.Next()
				return
			}

			// 前缀匹配（处理带参数的路径）
			if strings.HasPrefix(requestPath, path+"/") || strings.HasPrefix(path, requestPath+"/") {
				c.Next()
				return
			}
		}

		// 检查请求方法是否为禁止的类型（POST/PUT/DELETE）
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			response.Forbidden(c, "演示模式下禁止此操作")
			c.Abort()
			return
		}

		// 其他情况放行
		c.Next()
	}
}
