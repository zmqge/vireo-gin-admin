package models

import (
	"time"

	"gorm.io/gorm"
)

// NoticeReceiverModel NoticeReceiver实体
type NoticeReceiverModel struct {
	gorm.Model                // 自动包含ID字段（类型为uint）
	Name       string         `json:"name" gorm:"size:50;comment:NoticeReceiver名称"` // 名称字段
	CreatorID  uint           `json:"creator_id" gorm:"column:creator_id;index;comment:创建人ID"`
	DeptID     uint           `json:"dept_id" gorm:"column:dept_id;index;comment:部门ID"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (NoticeReceiverModel) TableName() string {
	return "notice_receiver" // 返回您想要的表名
}
