package models

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NoticesModel Notices实体
type NoticesModel struct {
	ID            uint                `json:"id" gorm:"primaryKey;autoIncrement"`
	Title         string              `json:"title" gorm:"size:50;comment:Notices名称"`
	Content       string              `json:"content" gorm:"size:255;comment:Notices内容"`
	Type          string              `json:"type" gorm:"default:0;comment:类型"`
	Level         string              `json:"level" gorm:"default:0;comment:级别"`
	TargetType    uint                `json:"targetType" gorm:"default:0;comment:'1:全体用户 2:指定部门 3:指定角色 4:指定用户'"`
	TargetIDs     []uint              `json:"targetIds" gorm:"column:target_ids;serializer:json;comment:目标用户ID"` // 改为uint数组
	Status        int                 `json:"publishStatus" gorm:"default:0;comment:状态"`
	IsRead        int                 `json:"isRead" gorm:"default:0;comment:是否已读"`
	PublisherName string              `json:"publisherName" gorm:"default:NULL;size:50;comment:发布人"`
	PublishedAt   time.Time           `json:"publishTime" gorm:"default:NULL;comment:发布时间"`
	RevokedAt     time.Time           `json:"revokeTime" gorm:"default:NULL;comment:撤回时间"`
	CreatorID     uint                `json:"creator_id" gorm:"column:creator_id;index;comment:创建人ID"`
	DeptID        uint                `json:"dept_id" gorm:"column:dept_id;index;comment:部门ID"`
	CreatedAt     time.Time           `json:"createTime" gorm:"comment:创建时间"`
	UpdatedAt     time.Time           `json:"-"`
	DeletedAt     gorm.DeletedAt      `json:"-" gorm:"index"`
	Receivers     []NoticeReceiverDTO `gorm:"-"`
}

// TableName 指定表名
func (NoticesModel) TableName() string {
	return "notices" // 返回您想要的表名
}

// NoticeReceiverDTO 接收者数据传输对象
type NoticeReceiverDTO struct {
	UserID uint `json:"user_id" binding:"required"`
}

// NoticeReceiver 接收者实体
type NoticeReceiver struct {
	NoticeID  uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"primaryKey"`
	IsRead    uint      `gorm:"default:0"`
	ReadTime  time.Time `gorm:"default:NULL"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (NoticeReceiver) TableName() string {
	return "notice_receiver" // 默认表名
}

// BeforeCreate 钩子函数，在创建前设置创建人ID和部门IDc
func (c *NoticesModel) BeforeCreate(db *gorm.DB) error {
	type User struct {
		ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
		Username string `json:"username" gorm:"size:100;not null;uniqueIndex;comment:用户名"`
		DeptID   int    `json:"dept_id" gorm:"column:dept_id;index;comment:部门ID"`
	}

	// 1. 从GORM上下文中获取gin.Context
	ctx, ok := db.Statement.Context.Value("ginContext").(*gin.Context)
	if !ok {
		log.Println("[NoticesModel] 错误：无法获取gin.Context")
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
				log.Printf("[NoticesModel] 转换 userID 为 uint 类型失败: %v", err)
			}
		}
	}
	c.CreatorID = UserID
	log.Printf("[NoticesModel] 已设置 CreatorID: %v", UserID)

	// 4. 查询部门ID
	var deptID uint
	if err := db.Model(User{}).
		Where("id = ?", UserID).
		Select("dept_id").
		First(&deptID).Error; err != nil {
		log.Printf("[NoticesModel] 查询部门ID失败: %v", err)
		return nil
	}

	c.DeptID = deptID
	return nil
}
