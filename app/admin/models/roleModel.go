// role.go
package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Role struct {
	ID        uint      `gorm:"primaryKey;comment:主键"`
	Name      string    `gorm:"size:50;comment:菜单名称"`
	Code      string    `gorm:"size:50;comment:菜单代码"`
	DataScope int       `gorm:"comment:数据范围"` // 数据范围 (1=全部数据, 2=自定义数据, 3=本部门及以下数据, 4=本部门数据, 5=仅本人数据)
	Sort      int       `gorm:"comment:排序"`
	Status    int       `gorm:"comment:是否隐藏"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	// DeletedAt   gorm.DeletedAt `gorm:"index;comment:删除时间"`   // 软删除字段，用于记录删除时间
	Users           []User       ` gorm:"many2many:user_roles;"`      // 用户列表（关联用户角色）
	Permissions     []Permission `gorm:"many2many:role_permissions;"` // 权限列表（关联角色权限）
	DeptIDs         []byte       `gorm:"type:json;comment:'自定义部门ID列表'"`
	PermissionDepts []uint       `gorm:"-"` // 临时存储的权限部门ID列表（不存储到数据库）
	CustomRoles     []*Role      `gorm:"-"` // 其他具有自定义权限的角色（临时存储）
}

// 用户-角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

// 角色-菜单关联表
type RoleMenu struct {
	RoleID uint `gorm:"primaryKey"`
	MenuID uint `gorm:"primaryKey"`
}
type RoleV0 struct {
	ID          uint   `json:"id" `         // 主键ID
	Name        string `json:"name"`        // 角色名称，唯一且非空
	Code        string `json:"code"`        // 角色代码，唯一且非空
	Status      int    `json:"status" `     // 状态 (1=启用, 0=禁用)
	Sort        int    `json:"sort" `       // 排序字段，默认为1
	DataScope   int    `json:"dataScope"`   // 数据范围 (1=全部数据, 2=自定义数据, 3=本部门及以下数据, 4=本部门数据, 5=仅本人数据)
	Description string `json:"description"` // 描述信息

}

// SetCustomDepts 设置角色的自定义部门ID列表
func (r *Role) SetCustomDepts(deptIDs []uint) error {
	if deptIDs == nil {
		r.DeptIDs = []byte("[]") // 空数组
		return nil
	}

	data, err := json.Marshal(deptIDs)
	if err != nil {
		return fmt.Errorf("部门ID序列化失败: %w", err)
	}

	r.DeptIDs = data
	return nil
}

// GetCustomDepts 获取角色的自定义部门ID列表
func (r *Role) GetCustomDepts() ([]uint, error) {
	if r.DeptIDs == nil || len(r.DeptIDs) == 0 {
		return []uint{}, nil
	}

	// 尝试处理可能的字符串格式(如逗号分隔的数字)
	strValue := string(r.DeptIDs)
	if strings.Contains(strValue, ",") {
		var deptIDs []uint
		parts := strings.Split(strings.Trim(strValue, "[]\" "), ",")
		for _, p := range parts {
			if id, err := strconv.ParseUint(strings.TrimSpace(p), 10, 32); err == nil {
				deptIDs = append(deptIDs, uint(id))
			}
		}
		return deptIDs, nil
	}

	// 尝试标准JSON解析
	var deptIDs []uint
	if err := json.Unmarshal(r.DeptIDs, &deptIDs); err != nil {
		// 如果解析失败，尝试解析单个数字
		if id, err := strconv.ParseUint(strValue, 10, 32); err == nil {
			return []uint{uint(id)}, nil
		}
		return nil, fmt.Errorf("解析部门ID列表失败: %w", err)
	}

	return deptIDs, nil
}
