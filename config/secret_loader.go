package config

import (
	"log"
	"os"
)

func validateSecret(secret string) bool {
	return len(secret) >= 32 // 确保密钥长度≥32字符
}

func LoadSecrets() {
	// 优先级：环境变量 > 配置文件
	if envSecret := os.Getenv("JWT_ACCESS_SECRET"); envSecret != "" {
		if !validateSecret(envSecret) {
			log.Fatal("JWT_ACCESS_SECRET 长度不足32字符")
		}
		App.JWT.AccessSecret = envSecret
		log.Println("[Config] JWT_ACCESS_SECRET loaded from env")
	}

	if envSecret := os.Getenv("JWT_REFRESH_SECRET"); envSecret != "" {
		App.JWT.RefreshSecret = envSecret
		log.Println("[Config] JWT_REFRESH_SECRET loaded from env")
	}
}
