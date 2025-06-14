// permission.go
package models

import (
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID          uint   `json:"id" gorm:"primaryKey"`        // 权限 ID
	Code        string `json:"code" gorm:"size:255"`        // 权限代码
	Name        string `json:"name" gorm:"size:255"`        // 权限名称
	Description string `json:"description" gorm:"size:255"` // 权限描述
	Type        string `json:"type" gorm:"size:50"`         // 权限类型（如 menu）
	Icon        string `json:"icon" gorm:"size:50"`         // 图标
	Module      string `json:"module" gorm:"size:50"`       // 模块（如 admin, user 等）
	ParentID    *uint  `json:"parent_id" gorm:"index"`      // 父权限 ID
}

type RolePermission struct {
	gorm.Model
	RoleID       uint `gorm:"index"`
	PermissionID uint `gorm:"index"`
}

// 用户-权限关联表
type UserPermission struct {
	gorm.Model
	UserID       uint `gorm:"not null"`
	PermissionID uint `gorm:"not null"`
}

func BuildPermissionTree(parentID uint) []map[string]interface{} {
	var permissions []Permission
	database.DB.Where("parent_id = ?", parentID).Find(&permissions)

	tree := make([]map[string]interface{}, 0)
	for _, p := range permissions {
		node := map[string]interface{}{
			"code":     p.Code,
			"name":     p.Name,
			"children": BuildPermissionTree(p.ID), // 递归构建子树
		}
		tree = append(tree, node)
	}
	return tree
}

func CheckUserPermission(userID uint, permissionCode string) bool {
	var count int64
	database.DB.Model(&UserRole{}).
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("user_roles.user_id = ? AND permissions.code = ?", userID, permissionCode).
		Count(&count)
	return count > 0

	// 方法2：直接用户-权限关联（如果启用）
	// db.Model(&UserPermission{}).
	//     Joins("JOIN permissions ON user_permissions.permission_id = permissions.id").
	//     Where("user_permissions.user_id = ? AND permissions.code = ?", userID, permissionCode).
	//     Count(&count)
}

func GetUserPermissionsFromDB(userID uint) []string {
	var permissions []string
	database.DB.Model(&UserRole{}).
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.code", &permissions)
	return permissions
}

type PermissionNode struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	Children []PermissionNode `json:"children,omitempty"`
}

type OptionPermLong struct {
	Value    string           `json:"value"`
	Label    string           `json:"label"`
	Children []OptionPermLong `json:"children,omitempty"`
}
