package redis_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
)

func TestMain(m *testing.M) {
	os.Setenv("PROJECT_ROOT", "E:/projects/vireo-gin-admin")
	config.Init()
	fmt.Printf("Redis 连接配置: Host=%s, Port=%s, DB=%d\n",
		config.App.Redis.Host,
		config.App.Redis.Port,
		config.App.Redis.DB,
	)

	// 初始化服务
	database.InitDB()
	redis.InitRedis()

	// 测试 Redis 连接
	if redis.Client == nil {
		panic("Redis 客户端未初始化")
	}
	if _, err := redis.Client.Ping(context.Background()).Result(); err != nil {
		panic("Redis 连接失败: " + err.Error())
	}

	code := m.Run()
	database.Close()
	os.Exit(code)
}

func TestBlacklist(t *testing.T) {
	testToken := "test_token_123"
	ctx := context.Background()
	key := "blacklist:" + testToken

	// 1. 添加 Token
	err := redis.AddToBlacklist(testToken, 10*time.Minute)
	if err != nil {
		t.Fatalf("添加黑名单失败: %v", err)
	}

	// 2. 打印 Key 详情
	val, err := redis.Client.Get(ctx, key).Result()
	fmt.Printf("Key: %s, Value: %s, Error: %v\n", key, val, err)

	// 3. 检查是否存在
	exists, err := redis.Client.Exists(ctx, key).Result()
	if err != nil {
		t.Fatalf("检查黑名单失败: %v", err)
	}
	fmt.Printf("Key 存在: %v\n", exists == 1)

	// 4. 断言
	if exists != 1 {
		t.Error("Token 应存在于黑名单中")
	}

	// 清理
	redis.Client.Del(ctx, key)
}

// 测试空 Token
func TestEmptyToken(t *testing.T) {
	err := redis.AddToBlacklist("", 10*time.Minute)
	if err == nil || err.Error() != "token 不能为空" {
		t.Errorf("空 Token 应返回错误，实际: %v", err)
	}
}

// 测试超短 TTL
func TestShortTTL(t *testing.T) {
	err := redis.AddToBlacklist("short_ttl_token", time.Millisecond)
	if err != nil {
		t.Fatalf("添加失败: %v", err)
	}
	time.Sleep(time.Millisecond * 2)
	exists, _ := redis.IsBlacklisted("short_ttl_token")
	if exists {
		t.Error("Key 应已过期")
	}
}

func BenchmarkBlacklist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token := fmt.Sprintf("token_%d", i)
		_ = redis.AddToBlacklist(token, time.Minute)
	}
}
