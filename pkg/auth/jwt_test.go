package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	j := NewJWT()
	userID := "123"

	// 生成 Token
	token, err := j.GenerateAccessToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 解析 Token
	claims, err := j.ParseAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.Subject)
}

func TestTokenExpiration(t *testing.T) {
	j := NewJWT()
	_, _ = j.GenerateAccessToken("123") // 忽略未使用的token

	// 模拟过期 Token
	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 已过期
	}).SignedString(j.AccessSecret)

	_, err := j.ParseAccessToken(expiredToken)
	assert.True(t, j.IsTokenExpired(err))
}
