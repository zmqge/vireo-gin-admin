// user.go
package models

import (
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Status   int    `json:"status"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	OpenId   string `json:"open_id"` // 微信小程序或公众号的 OpenId

	DeptID   uint   `json:"dept_id"`
	RoleList []Role `json:"role_list" gorm:"many2many:user_roles;"`
	Dept     Dept   `gorm:"foreignKey:DeptID"`

	// 权限相关字段（不映射到数据库）
	DataScope       int    `json:"data_scope" gorm:"-"` // 用户最高数据权限（不存储）
	PermissionDepts []uint `json:"-" gorm:"-"`          // 权限部门ID列表（临时存储）
}

func (u *User) TableName() string {
	return "users"
}

type UserPageVO struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Mobile     string `json:"mobile"`
	Gender     string `json:"gender"`
	Avatar     string `json:"avatar"`
	Email      string `json:"email"`
	Status     int    `json:"status"`
	DeptID     uint   `json:"dept_id"`
	DeptName   string `json:"deptName"`   // 部门名称
	RoleNames  string `json:"roleNames"`  // 角色名称，逗号分隔
	CreateTime string `json:"createTime"` // 创建时间
}

type UserQueryParams struct {
	Keywords   string
	Status     string
	RoleIDs    []string // 使用切片代替逗号分隔的字符串
	CreateTime []string // 使用切片代替逗号分隔的字符串
	Field      string
	Direction  string
	PageNum    int
	PageSize   int
	DeptID     int
}

type UserPageResult struct {
	Users []User
	Total int64
}

type DataUserPageVO struct {
	List  []UserPageVO `json:"list"`
	Total int64        `json:"total"`
}
type UserOption struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// 统一的密码加密函数，自动拼接 password+salt
func HashPasswordWithSalt(password, salt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func IsUserSuperAdmin(userID uint) bool {
	var count int64
	database.DB.Model(&UserRole{}).
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, "super_admin").
		Count(&count)
	return count > 0
}
