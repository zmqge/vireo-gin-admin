package models

import (
	"encoding/json"
	"time"
)

type Menu struct {
	ID         uint      `gorm:"primaryKey;comment:主键"`
	ParentID   uint      `gorm:"column:parent_id;comment:父菜单ID"`
	Name       string    `gorm:"size:50;comment:菜单名称"`
	Path       string    `gorm:"size:100;comment:路由路径"`
	Component  string    `gorm:"size:100;comment:前端组件路径"`
	Redirect   string    `gorm:"size:100;comment:跳转链接"`
	Perm       string    `gorm:"size:100;comment:权限标识"`
	Icon       string    `gorm:"size:50;comment:图标"`
	Sort       int       `gorm:"comment:排序"`
	Visible    int       `gorm:"comment:是否隐藏"`
	Type       int       `gorm:"comment:菜单类型"`
	Title      string    `gorm:"size:50;comment:路由标题"`
	KeepAlive  int       `gorm:"comment:是否开启页面缓存"`
	AlwaysShow int       `gorm:"comment:是否始终显示"`
	Params     string    `gorm:"type:text;comment:路由参数"`
	CreatedAt  time.Time `gorm:"comment:创建时间"`
	UpdatedAt  time.Time `gorm:"comment:更新时间"`
	// DeletedAt  gorm.DeletedAt `gorm:"index;comment:删除时间"`     // 软删除字段，用于记录删除时间
	Permission []Permission `gorm:"many2many:menu_permissions;"`
}
type MenuVO struct {
	ID        uint      `json:"id"`
	ParentID  uint      `json:"parentID"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	RouteName string    `json:"routeName"`
	RoutePath string    `json:"routePath"`
	Component string    `json:"component"`
	Sort      int       `json:"sort"`
	Visible   int       `json:"visible"`
	Icon      string    `json:"icon"`
	Redirect  string    `json:"redirect"`
	Children  *[]MenuVO `json:"children,omitempty"` // 使用指针和omitempty
}
type MenuPermission struct {
	MenuID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

// JSON 将 map 转换为 JSON 字符串
func JSON(m map[string]interface{}) string {
	data, _ := json.Marshal(m)
	return string(data)
}

type RouteVO struct {
	Path      string     `json:"path"`
	Component string     `json:"component"`
	Redirect  string     `json:"redirect"`
	Name      string     `json:"name"`
	ParentID  uint       `json:"parentID"`
	Meta      Meta       `json:"meta"`
	Children  *[]RouteVO `json:"children,omitempty"` // 使用指针和omitempty
}

// Meta 路由属性类型
type Meta struct {
	Title      string `json:"title"`      // 路由title
	Icon       string `json:"icon"`       // ICON
	Hidden     bool   `json:"hidden"`     // 是否隐藏
	KeepAlive  bool   `json:"keepAlive"`  // 是否开启页面缓存
	AlwaysShow bool   `json:"alwaysShow"` // 是否始终显示
	Params     string `json:"params"`     // 路由参数
}

// 假设 OptionLong 是一个自定义结构体，这里添加其定义，实际使用时请根据项目情况调整
type OptionLong struct {
	Value    uint         `json:"value"`
	Label    string       `json:"label"`
	Children []OptionLong `json:"children,omitempty"`
}
