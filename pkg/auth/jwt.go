package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
)

// 移除models.User依赖，改用简化结构
type UserInfo struct {
	ID          uint
	Username    string
	RoleCodes   []string // 替换GetRoleCodes()
	Permissions []string // 替换GetPermissions()
}

type JWT struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

// 确保CustomClaims结构体存在
type CustomClaims struct {
	UserID      uint
	Username    string
	Roles       []string
	Permissions []string
	jwt.RegisteredClaims
}

// 初始化JWT工具
func NewJWT() *JWT {
	return &JWT{
		AccessSecret:  []byte(config.App.JWT.AccessSecret),
		RefreshSecret: []byte(config.App.JWT.RefreshSecret),
		AccessExpire:  config.App.JWT.AccessExpire,
		RefreshExpire: config.App.JWT.RefreshExpire,
	}
}

// 生成双Token
func (j *JWT) GenerateTokens(userID uint, username string) (string, string, error) {
	// Access Token (15分钟过期)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		Issuer:    "vireo-gin-admin",
	})
	accessSigned, err := accessToken.SignedString(j.AccessSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (7天有效期)
	refreshToken := uuid.New().String()
	if err := redis.Client.Set(context.Background(), fmt.Sprintf("refresh:%d", userID), refreshToken, 7*24*time.Hour).Err(); err != nil {
		return "", "", err
	}

	return accessSigned, refreshToken, nil
}

// 解析Access Token
func (j *JWT) ParseAccessToken(tokenString string) (*jwt.RegisteredClaims, error) {
	return parseToken(tokenString, j.AccessSecret)
}

// 解析Refresh Token
func (j *JWT) ParseRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	if j == nil || len(j.RefreshSecret) == 0 {
		return nil, errors.New("JWT实例未初始化或RefreshSecret为空")
	}
	return parseToken(tokenString, j.RefreshSecret)
}

func parseToken(tokenString string, secret []byte) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// 生成带权限信息的Token
func (j *JWT) GenerateToken(user UserInfo) (string, error) {
	claims := CustomClaims{
		UserID:      user.ID,
		Username:    user.Username,
		Roles:       user.RoleCodes,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessExpire)),
		},
	}
	return j.CreateToken(claims)
}

// ParseToken 解析 JWT Token
func (j *JWT) ParseToken(c *gin.Context) (*jwt.RegisteredClaims, error) {
	tokenString := extractToken(c)
	if tokenString == "" {
		return nil, jwt.ErrInvalidKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.AccessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.AccessSecret)
}

// 判断 Token 是否过期
func (j *JWT) IsTokenExpired(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired)
}

// 生成仅 Access Token
func (j *JWT) GenerateAccessToken(userID string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}).SignedString(j.AccessSecret)
}

// JWTAuthRefresh 支持自动续期的中间件
func (j *JWT) JWTAuthRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供 Token"})
			return
		}

		claims, err := j.ParseToken(c)
		if err != nil {
			if j.IsTokenExpired(err) {
				// 检查 Refresh Token 是否存在
				refreshToken, err := redis.Client.Get(
					context.Background(),
					fmt.Sprintf("refresh:%s", claims.Subject),
				).Result()
				if err != nil || refreshToken == "" {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh Token 无效"})
					return
				}

				// 生成新 Access Token
				newToken, err := j.GenerateAccessToken(claims.Subject)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
					return
				}

				// 返回新 Token
				c.Header("New-Access-Token", newToken)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token 无效"})
			return
		}

		c.Next()
	}
}

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
