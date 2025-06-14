package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmqge/vireo-gin-admin/config"
)

func TestLoadSecrets(t *testing.T) {
	// 备份原始配置
	originalAccess := config.App.JWT.AccessSecret
	originalRefresh := config.App.JWT.RefreshSecret

	t.Run("环境变量覆盖配置", func(t *testing.T) {
		t.Setenv("JWT_ACCESS_SECRET", "test_access")
		t.Setenv("JWT_REFRESH_SECRET", "test_refresh")

		config.LoadSecrets()

		assert.Equal(t, "test_access", config.App.JWT.AccessSecret)
		assert.Equal(t, "test_refresh", config.App.JWT.RefreshSecret)
	})

	// 恢复配置
	config.App.JWT.AccessSecret = originalAccess
	config.App.JWT.RefreshSecret = originalRefresh
}
