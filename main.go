package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"github.com/zmqge/vireo-gin-admin/pkg/middleware"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
	"github.com/zmqge/vireo-gin-admin/routes"
	"go.uber.org/zap"
)

func main() {
	// 设置配置文件路径和名称
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config") // 添加 config 目录作为配置文件搜索路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	// 获取需要扫描的目录路径
	controllerDirs := viper.GetStringSlice("controller_dirs")
	fmt.Println("Controller directories to scan:", controllerDirs)

	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// 1. 初始化配置
	config.Init()

	// 2. 初始化数据库
	db := database.InitDB()
	redis.InitRedis()
	defer database.Close()

	// 创建 Gin 引擎
	r := gin.Default()

	//注入数据库到上下文
	r.Use(func(c *gin.Context) {
		c.Set("db", db)              // 数据库
		c.Set("redis", redis.Client) // 如果需要Redis也注入
		c.Next()
	})

	// 应用全局中间件
	r.Use(middleware.Logger())   // 日志中间件
	r.Use(middleware.Recovery()) // 恢复中间件
	r.Use(middleware.Cors())     // 跨域中间件
	// 注册所有路由
	routes.RegisterAllRoutes(r, db)

	// 显示欢迎画面
	showWelcomeMessage()

	// 启动服务
	r.Run(":8080")
}

// showWelcomeMessage 显示欢迎画面
func showWelcomeMessage() {
	fmt.Printf(`
 __      __  _____   _____    ______    ____  
 \ \    / / |_   _| |  __ \  |  ____|  / __ \ 
  \ \  / /    | |   | |__) | | |__    | |  | |
   \ \/ /     | |   |  _  /  |  __|   | |  | |
    \  /     _| |_  | | \ \  | |____  | |__| |
     \/     |_____| |_|  \_\ |______|  \____/ 
                                              
                                              
============================================
  服务启动成功! 
  ➤ 时间: %s
  ➤ 环境: %s
  ➤ 地址: http://localhost:%s
  ➤ 文档: http://localhost:%s/swagger/index.html
============================================
`, time.Now().Format("2006-01-02 15:04:05"), gin.Mode(), "8080", "8080")
}
