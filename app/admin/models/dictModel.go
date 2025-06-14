package models

import (
	"gorm.io/gorm"
)

// DictModel Dict实体
type DictModel struct {
	gorm.Model        // 自动包含ID字段（类型为uint）
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:Dict主键ID"`
	Name       string `json:"name" gorm:"size:50;comment:Dict名称"`
	DictCode   string `json:"dictCode" gorm:"size:50;comment:Dict编码"`
	Status     int    `json:"status" gorm:"default:1;comment:Dict状态 1启用 0禁用"`
	Remark     string `json:"remark" gorm:"size:255;comment:Dict备注"`
	Sort       int    `json:"sort" gorm:"default:0;comment:Dict排序"`
}

// TableName 指定表名
func (DictModel) TableName() string {
	return "dict_type" // 返回您想要的表名
}

// DictItemModel DictItem实体
type DictItemModel struct {
	// gorm.Model        // 自动包含ID字段（类型为uint）
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:DictItem主键ID"`
	DictCode string `json:"dictCode" gorm:"size:50;not null;comment:DictItem编码"`
	Value    string `json:"value" gorm:"size:50;not null;comment:DictItem值"`
	Label    string `json:"label" gorm:"size:50;not null;comment:DictItem标签"`
	TagType  string `json:"tagType" gorm:"size:50;not null;comment:DictItem标签类型"`
	Status   int    `json:"status" gorm:"default:1;comment:DictItem状态 1启用 0禁用"`
	Sort     int    `json:"sort" gorm:"default:0;comment:DictItem排序"`
}

// TableName 指定表名
func (DictItemModel) TableName() string {
	return "dict_item" // 返回您想要的表名
}
