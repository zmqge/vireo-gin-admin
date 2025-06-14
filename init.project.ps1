# 定义项目根目录
$projectRoot = "E:\projects\vireo-gin-admin"

# 创建文件夹结构
$directories = @(
    "$projectRoot\app\controllers",
    "$projectRoot\app\middleware",
    "$projectRoot\app\models",
    "$projectRoot\app\services",
    "$projectRoot\app\validators",
    "$projectRoot\config",
    "$projectRoot\database",
    "$projectRoot\docs",
    "$projectRoot\pkg\auth",
    "$projectRoot\pkg\constant",
    "$projectRoot\pkg\errors",
    "$projectRoot\pkg\logger",
    "$projectRoot\pkg\response",
    "$projectRoot\routes",
    "$projectRoot\scripts",
    "$projectRoot\storage",
    "$projectRoot\utils"
)

foreach ($dir in $directories) {
    if (!(Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force
    }
}

# 创建文件并写入内容
$files = @{
    "$projectRoot\config\config.go" = @'
package config

import "github.com/spf13/viper"

type AppConfig struct {
    Port     string `mapstructure:"PORT"`
    Database struct {
        Host     string `mapstructure:"HOST"`
        Port     string `mapstructure:"PORT"`
        User     string `mapstructure:"USER"`
        Password string `mapstructure:"PASSWORD"`
        Name     string `mapstructure:"NAME"`
    } `mapstructure:"DATABASE"`
    JWT struct {
        SignKey    string `mapstructure:"SIGN_KEY"`
        ExpireTime int    `mapstructure:"EXPIRE_TIME"`
    } `mapstructure:"JWT"`
}

var App AppConfig

func Init() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    if err := viper.ReadInConfig(); err != nil {
        panic("配置文件读取失败: " + err.Error())
    }

    if err := viper.Unmarshal(&App); err != nil {
        panic("配置文件解析失败: " + err.Error())
    }
}
'@

    "$projectRoot\config\config.yaml" = @'
APP:
  PORT: "8080"

DATABASE:
  HOST: "localhost"
  PORT: "3306"
  USER: "root"
  PASSWORD: "password"
  NAME: "vireo_gin_admin"

JWT:
  SIGN_KEY: "your-secret-key"
  EXPIRE_TIME: 3600
'@

    "$projectRoot\pkg\logger\logger.go" = @'
package logger

import "log"

func Init() {
    log.Println("Logger initialized")
}

func Fatal(format string, v ...interface{}) {
    log.Fatalf(format, v...)
}
'@

    "$projectRoot\routes\admin.go" = @'
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/zmqge/vireo-gin-admin/app/controllers"
)

func RegisterAdminRoutes(r *gin.Engine) {
    ctrl := controllers.AdminController{}
    admin := r.Group("/admin")
    {
        admin.GET("/ping", ctrl.Ping)
    }
}
'@

    "$projectRoot\app\controllers\admin.go" = @'
package controllers

import "github.com/gin-gonic/gin"

type AdminController struct{}

func (ctrl *AdminController) Ping(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
}
'@

    "$projectRoot\main.go" = @'
package main

import (
    "github.com/zmqge/vireo-gin-admin/config"
    "github.com/zmqge/vireo-gin-admin/pkg/logger"
    "github.com/zmqge/vireo-gin-admin/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化配置
    config.Init()

    // 初始化日志
    logger.Init()

    // 创建Gin引擎
    r := gin.Default()

    // 注册路由
    routes.RegisterAdminRoutes(r)

    // 启动服务
    if err := r.Run(":" + config.App.Port); err != nil {
        logger.Fatal("服务启动失败: %v", err)
    }
}
'@

    "$projectRoot\go.mod" = @'
module github.com/zmqge/vireo-gin-admin

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/spf13/viper v1.16.0
    gorm.io/gorm v1.25.2
    gorm.io/driver/mysql v1.5.1
)
'@
}

foreach ($file in $files.GetEnumerator()) {
    Set-Content -Path $file.Key -Value $file.Value -Force
}

Write-Host "项目初始化完成！请运行以下命令启动项目："
Write-Host "cd $projectRoot"
Write-Host "go mod tidy"
Write-Host "go run main.go"
