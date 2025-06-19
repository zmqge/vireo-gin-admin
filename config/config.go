package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string `mapstructure:"PORT"`
	Database struct {
		Host     string `mapstructure:"HOST"`
		Port     string `mapstructure:"PORT"`
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
		Name     string `mapstructure:"NAME"`
	} `mapstructure:"DATABASE"`
	Redis struct {
		Host     string        `mapstructure:"HOST"`
		Port     string        `mapstructure:"PORT"`
		Password string        `mapstructure:"PASSWORD"`
		DB       int           `mapstructure:"DB"`
		TTL      time.Duration `mapstructure:"TTL"`
	} `mapstructure:"REDIS"`
	JWT struct {
		Secret        string        `mapstructure:"SECRET"`
		ExpireTime    time.Duration `mapstructure:"EXPIRE_TIME"`
		AccessSecret  string        `mapstructure:"ACCESS_SECRET"`
		RefreshSecret string        `mapstructure:"REFRESH_SECRET"`
		AccessExpire  time.Duration `mapstructure:"ACCESS_EXPIRE"`
		RefreshExpire time.Duration `mapstructure:"REFRESH_EXPIRE"`
	} `mapstructure:"JWT"`
	RBAC struct {
		CacheTTL       int    `mapstructure:"CacheTTL"`
		SuperAdminRole string `mapstructure:"SuperAdminRole"`
		AdminRole      string `mapstructure:"AdminRole"`
	} `mapstructure:"RBAC"`
	ControllerDirs []string `mapstructure:"CONTROLLER_DIRS"`
	DemoMode       bool     `mapstructure:"DEMO_MODE"`
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
}

var App Config

func Init() {
	// 1. 先加载配置文件
	loadYamlConfig()

	// 2. 用环境变量覆盖配置
	LoadSecrets()

	// 3. 验证密钥非空
	if App.JWT.AccessSecret == "" || App.JWT.RefreshSecret == "" {
		panic("JWT 密钥不能为空")
	}
}

func getProjectRoot() string {
	// 方法1：通过环境变量指定（推荐）
	if root := os.Getenv("PROJECT_ROOT"); root != "" {
		return root
	}

	// 方法2：向上递归查找 go.mod
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // 到达根目录
		}
		dir = parent
	}
	panic("未找到项目根目录（需包含 go.mod）")
}

func loadYamlConfig() {
	root := getProjectRoot()
	configPath := filepath.Join(root, "config", "config.yaml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}

	// 打印所有加载的配置
	fmt.Println("全部配置:", viper.AllSettings())

	if err := viper.Unmarshal(&App); err != nil {
		panic("配置解析失败: " + err.Error())
	}

	// 显式检查 Port 字段
	fmt.Printf("App.Port: %#v\n", App.Port)
}
