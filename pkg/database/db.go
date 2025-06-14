package database

import (
	"fmt"

	"github.com/zmqge/vireo-gin-admin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 初始化数据库连接
func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.App.Database.User,
		config.App.Database.Password,
		config.App.Database.Host,
		config.App.Database.Port,
		config.App.Database.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	return DB
}

// 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// 新增 Close 方法
func Close() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
	}
}
