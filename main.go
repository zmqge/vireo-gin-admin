package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"github.com/zmqge/vireo-gin-admin/pkg/middleware"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
	"github.com/zmqge/vireo-gin-admin/routes"
	"go.uber.org/zap"
)

func main() {

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
	log.Println("AllowedOrigins:", config.App.AllowedOrigins)
	// r.Use(middleware.Cors())     // 跨域中间件
	r.Use(cors.New(cors.Config{
		// 允许所有来源，生产环境建议指定具体域名
		AllowOrigins: config.App.AllowedOrigins,
		// 允许的 HTTP 方法
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 允许的请求头
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// 是否允许携带凭证（如 Cookie）
		AllowCredentials: true,
		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	}))
	r.Use(middleware.DemoMode())
	r.Use(middleware.Logger())   // 日志中间件
	r.Use(middleware.Recovery()) // 恢复中间件

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
