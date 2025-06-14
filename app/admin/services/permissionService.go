package services

import (
	"fmt"

	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/annotations"
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"gorm.io/gorm"
)

type PermissionService struct {
	DB *gorm.DB
}

func NewPermissionService() *PermissionService {
	return &PermissionService{
		DB: database.DB, // 确保 database.DB 已正确初始化
	}
}

func (s *PermissionService) CheckPermission(userID, permission string) bool {
	// 实际权限检查逻辑
	var count int64
	s.DB.Model(&models.UserRole{}).
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("user_roles.user_id = ? AND permissions.name = ?", userID, permission).
		Count(&count)

	return count > 0
}

// 定义权限树节点结构
type PermissionNode struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Children    []PermissionNode `json:"children"`
}

func (s *PermissionService) GetPermissionTree() ([]models.PermissionNode, error) {
	// 实现权限树逻辑
	return nil, nil
}

// SyncPermissions 同步权限
func (s *PermissionService) SyncPermissions() error {
	// 扫描多个目录
	permissions, err := annotations.ScanDirectories()
	if err != nil {
		return fmt.Errorf("扫描目录失败: %v", err)
	}

	// 将权限同步到数据库
	for _, perm := range permissions {
		// 处理单权限
		if err := database.DB.Create(&models.Permission{
			Code:        perm.Code,
			Name:        perm.Name,
			Description: perm.Description,
		}).Error; err != nil {
			return fmt.Errorf("插入权限失败: %v", err)
		}
	}
	return nil
}

// GetUserRolesAndPermissions 获取用户的角色和权限
func (s *PermissionService) GetUserRolesAndPermissions(userID string) ([]string, []string, error) {
	// 记录输入的用户 ID
	fmt.Printf("开始查询用户角色和权限，用户 ID: %s\n", userID)

	// 查询用户的角色
	var roles []string
	err := s.DB.Model(&models.UserRole{}).
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("roles.name", &roles).Error
	if err != nil {
		// 记录查询角色失败的错误
		fmt.Printf("查询用户角色失败，用户 ID: %s, 错误: %v\n", userID, err)
		return nil, nil, fmt.Errorf("查询用户角色失败: %v", err)
	}

	// 记录查询到的角色
	fmt.Printf("用户角色查询成功，用户 ID: %s, 角色: %v\n", userID, roles)

	// 查询用户的权限
	var permissions []string
	err = s.DB.Model(&models.UserRole{}).
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("role_permissions.permission_code", &permissions).Error
	if err != nil {
		// 记录查询权限失败的错误
		fmt.Printf("查询用户权限失败，用户 ID: %s, 错误: %v\n", userID, err)
		return nil, nil, fmt.Errorf("查询用户权限失败: %v", err)
	}

	// 记录查询到的权限
	fmt.Printf("用户权限查询成功，用户 ID: %s, 权限: %v\n", userID, permissions)

	return roles, permissions, nil
}

// ListPermOptions 列出权限选项（按模块分组，OptionLong格式）
func (s *PermissionService) ListPermOptions() ([]models.OptionPermLong, error) {
	var permissions []models.Permission
	if err := s.DB.Order("module ASC, id ASC").Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("查询权限选项失败: %v", err)
	}
	// 按 module 分组
	groupMap := make(map[string][]models.OptionPermLong)
	for _, perm := range permissions {
		opt := models.OptionPermLong{
			Value: perm.Code, // OptionPermLong.Value 类型应与定义一致
			Label: perm.Name + "   【" + perm.Code + "】",
		}
		groupMap[perm.Module] = append(groupMap[perm.Module], opt)
	}
	// 组装分组 OptionPermLong
	var result []models.OptionPermLong
	for module, opts := range groupMap {
		result = append(result, models.OptionPermLong{
			Value:    "0", // 分组节点 value 设为0
			Label:    module,
			Children: opts,
		})
	}
	return result, nil
}
