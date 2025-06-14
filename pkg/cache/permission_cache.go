package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
)

const PermissionCacheKey = "user_permissions:%d" // 用户ID占位符

// 获取用户权限列表（优先从缓存读取）
func GetUserPermissions(userID uint) ([]string, error) {
	ctx := context.Background()
	key := fmt.Sprintf(PermissionCacheKey, userID)

	// 从Redis获取
	if permissions, err := redis.Client.SMembers(ctx, key).Result(); err == nil && len(permissions) > 0 {
		return permissions, nil
	}

	// 从数据库加载
	permissions := models.GetUserPermissionsFromDB(userID)
	if _, err := redis.Client.SAdd(ctx, key, permissions).Result(); err != nil {
		return nil, err
	}
	redis.Client.Expire(ctx, key, 24*time.Hour)
	return permissions, nil
}

// 清除用户权限缓存
func ClearUserPermissionsCache(userID uint) {
	ctx := context.Background()
	redis.Client.Del(ctx, fmt.Sprintf(PermissionCacheKey, userID))
}
