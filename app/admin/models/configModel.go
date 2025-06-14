package models

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfigModel Config实体
type ConfigModel struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ConfigName  string         `json:"configName" gorm:"size:100;comment:配置名称"`
	ConfigKey   string         `json:"configKey" gorm:"size:100;not null;uniqueIndex;comment:配置键"`
	ConfigValue string         `json:"configValue" gorm:"size:500;comment:配置值"`
	Remark      string         `json:"remark" gorm:"size:500;comment:描述备注"`
	CreatorID   uint           `json:"creator_id" gorm:"column:creator_id;index;comment:创建人ID"`
	DeptID      uint           `json:"dept_id" gorm:"column:dept_id;index;comment:部门ID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (ConfigModel) TableName() string {
	return "Config" // 返回您想要的表名
}

// BeforeCreate 钩子函数，在创建前设置创建人ID和部门IDc
func (c *ConfigModel) BeforeCreate(db *gorm.DB) error {
	type User struct {
		ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
		Username string `json:"username" gorm:"size:100;not null;uniqueIndex;comment:用户名"`
		DeptID   int    `json:"dept_id" gorm:"column:dept_id;index;comment:部门ID"`
	}

	// 1. 从GORM上下文中获取gin.Context
	ctx, ok := db.Statement.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("[ConfigModel] 错误：无法获取gin.Context")
		return nil
	}
	// 将从上下文获取的 userID 字符串转换为 uint 类型
	var UserID uint
	if idStr, ok := ctx.Get("userID"); ok {
		if idStrStr, ok := idStr.(string); ok {
			var err error
			parsedID, err := strconv.ParseUint(idStrStr, 10, 64)
			UserID = uint(parsedID)
			if err != nil {
				log.Printf("[ConfigModel] 转换 userID 为 uint 类型失败: %v", err)
			}
		}
	}
	c.CreatorID = UserID
	log.Printf("[ConfigModel] 已设置 CreatorID: %v", UserID)

	// 4. 查询部门ID
	var deptID uint
	if err := db.Model(User{}).
		Where("id = ?", UserID).
		Select("dept_id").
		First(&deptID).Error; err != nil {
		log.Printf("[ConfigModel] 查询部门ID失败: %v", err)
		return nil
	}

	c.DeptID = deptID
	return nil
}
