package database

import (
	"context"

	"github.com/zmqge/vireo-gin-admin/pkg/redis"
	"gorm.io/gorm"
)

// 监听权限变更
func RegisterHooks(db *gorm.DB) {
	// 角色权限变更时清除缓存
	db.Callback().Update().After("gorm:update").Register("clear_permission_cache", func(tx *gorm.DB) {
		if tx.Statement.Table == "role_permissions" {
			redis.Client.Del(context.Background(), "user:*:perms")
		}
	})
}
