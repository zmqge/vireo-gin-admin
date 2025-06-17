// middleware/demo_mode.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

func DemoMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.App.DemoMode && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE") {
			response.Forbidden(c, "演示模式下禁止此操作")
			c.Abort()
			return
		}
	}
}
