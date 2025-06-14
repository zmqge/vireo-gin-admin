package services

import (
	"context"
	"fmt"

	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

// 存储 Refresh Token 到 Redis
func (s *TokenService) SaveRefreshToken(userID uint, token string) error {
	key := fmt.Sprintf("refresh_token:%d", userID)
	ctx := context.Background()
	return redis.Client.Set(ctx, key, token, config.App.Redis.TTL).Err()
}

// 从 Redis 获取 Refresh Token
func (s *TokenService) GetRefreshToken(userID uint) (string, error) {
	key := fmt.Sprintf("refresh_token:%d", userID)
	ctx := context.Background()
	return redis.Client.Get(ctx, key).Result()
}

// 从 Redis 删除 Refresh Token
func (s *TokenService) DeleteRefreshToken(userID int) error {
	key := fmt.Sprintf("refresh_token:%d", userID)
	ctx := context.Background()
	return redis.Client.Del(ctx, key).Err()
}
