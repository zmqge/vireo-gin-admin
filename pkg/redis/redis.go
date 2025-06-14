package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zmqge/vireo-gin-admin/config"
)

var Client *redis.Client

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.App.Redis.Host + ":" + config.App.Redis.Port,
		Password: config.App.Redis.Password,
		DB:       config.App.Redis.DB,
	})

	// 测试连接
	if _, err := Client.Ping(context.Background()).Result(); err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
}

func Init() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.App.Redis.Host, config.App.Redis.Port),
		Password: config.App.Redis.Password,
		DB:       config.App.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	return err
}

func AddToBlacklist(token string, expiration time.Duration) error {
	if token == "" {
		return errors.New("token 不能为空")
	}
	key := "blacklist:" + token
	return Client.Set(context.Background(), key, "1", expiration).Err()
}

func IsBlacklisted(token string) (bool, error) {
	if token == "" {
		return false, errors.New("token 不能为空")
	}
	key := "blacklist:" + token
	exists, err := Client.Exists(context.Background(), key).Result()
	return exists == 1, err
}
