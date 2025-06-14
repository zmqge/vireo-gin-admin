package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zmqge/vireo-gin-admin/pkg/auth"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

// Claims JWT 声明结构体
type Claims struct {
	UserID               string   `json:"user_id"`
	Username             string   `json:"username"`
	Roles                []string `json:"roles"`
	jwt.RegisteredClaims          // 嵌入 jwt.RegisteredClaims 以实现 jwt.Claims 接口
}

// JWT JWT 中间件
func JWT() gin.HandlerFunc {
	jwt := auth.NewJWT()
	return func(c *gin.Context) {
		// 1. 从请求头中提取 Token
		tokenString := extractToken(c)
		if tokenString == "" {
			response.Unauthorized(c, "Token缺失")
			c.Abort()
			return
		}

		// 2. 调用 auth/jwt.go 中的 ParseAccessToken 方法解析 Token
		claims, err := jwt.ParseAccessToken(tokenString)
		if err != nil {
			if jwt.IsTokenExpired(err) {
				response.Unauthorized(c, "Token已过期")
			} else {
				response.Unauthorized(c, "Token无效")
			}
			c.Abort()
			return
		}

		// 3. 将用户信息存储到上下文中
		c.Set("userID", claims.Subject)                                                // 或 claims.UserID
		log.Printf("[JWT 调试] 已设置 userID: %v (类型: %T)", claims.Subject, claims.Subject) // 关键调试

		// 立即验证是否能读取
		if val, exists := c.Get("userID"); exists {
			log.Printf("[JWT 调试] 上下文验证成功: userID=%v", val)
		} else {
			log.Printf("[JWT 调试] 错误: 上下文设置失败！")
		}

		c.Next()
	}
}

// extractToken 从请求头中提取 Token
func extractToken(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return ""
	}

	// 确保格式为 "Bearer <token>"
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
