package services

import (
	"context"
	"fmt"
	"time"

	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
)

// TokenService 定义令牌服务接口
type TokenService interface {
	SaveRefreshToken(userID uint, token, clientIP string, loginTime time.Time) error
	GetRefreshToken(userID uint) (string, error)
	GetLoginInfo(userID uint) (map[string]string, error)
	DeleteRefreshToken(userID uint) error // 修复参数类型与接口一致
}

// TokenServiceImpl 实现 TokenService 接口
type TokenServiceImpl struct{}

// NewTokenService 创建 TokenService 实例
func NewTokenService() TokenService { // 修复返回类型
	return &TokenServiceImpl{}
}

// 存储 Refresh Token 到 Redis
func (s *TokenServiceImpl) SaveRefreshToken(userID uint, token, clientIP string, loginTime time.Time) error {
	ctx := context.Background()
	tokenKey := fmt.Sprintf("refresh_token:%d", userID)
	loginInfoKey := fmt.Sprintf("login_info:%d", userID)

	// 存储刷新令牌
	err := redis.Client.Set(ctx, tokenKey, token, config.App.Redis.TTL).Err()
	if err != nil {
		return err
	}

	// 存储登录信息（单独的 key，与令牌共享相同 TTL）
	loginInfo := map[string]interface{}{
		"client_ip":  clientIP,
		"login_time": loginTime.Format(time.RFC3339),
	}

	err = redis.Client.HMSet(ctx, loginInfoKey, loginInfo).Err()
	if err != nil {
		return err
	}

	// 为登录信息设置相同的过期时间
	return redis.Client.Expire(ctx, loginInfoKey, config.App.Redis.TTL).Err()
}

// 从 Redis 获取 Refresh Token
func (s *TokenServiceImpl) GetRefreshToken(userID uint) (string, error) {
	key := fmt.Sprintf("refresh_token:%d", userID)
	ctx := context.Background()

	return redis.Client.Get(ctx, key).Result()
}

// 获取用户登录信息
func (s *TokenServiceImpl) GetLoginInfo(userID uint) (map[string]string, error) {
	key := fmt.Sprintf("login_info:%d", userID)
	ctx := context.Background()

	return redis.Client.HGetAll(ctx, key).Result()
}

// 从 Redis 删除 Refresh Token
func (s *TokenServiceImpl) DeleteRefreshToken(userID uint) error { // 修复参数类型
	key := fmt.Sprintf("refresh_token:%d", userID)
	ctx := context.Background()
	return redis.Client.Del(ctx, key).Err()
}
