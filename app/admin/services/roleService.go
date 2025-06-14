package services

import (
	"errors"
	"fmt"

	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
	"gorm.io/gorm"
)

type RoleService struct {
	repo *repositories.RoleRepository
	db   *gorm.DB // 事务等特殊场景可用
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		repo: repositories.NewRoleRepository(db),
		db:   db,
	}
}

// CreateRole 创建新角色
func (s *RoleService) CreateRole(role *models.Role) error {
	// 先检查角色名称/编码唯一性
	var roles []models.Role
	roles, _ = s.repo.ListRoles(role.Name, "", "", 1, 1)
	for _, r := range roles {
		if r.Name == role.Name || r.Code == role.Code {
			return fmt.Errorf("角色名称 '%s' 或者角色编码 '%s'已存在", role.Name, role.Code)
		}
	}
	return s.repo.CreateRole(role)
}

// GetRolePage 获取角色分页列表
func (s *RoleService) GetRolePage(keywords, startDate, endDate string, pageNum, pageSize int) (interface{}, error) {
	roles, err := s.repo.ListRoles(keywords, startDate, endDate, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	vos := RolesToV0List(roles)
	result := map[string]interface{}{
		"list":  vos,
		"total": len(roles), // 如需精确total可扩展repo返回
	}
	return result, nil
}

// GetRoleDetails 获取角色详情
func (s *RoleService) GetRoleDetails(id string) (interface{}, error) {
	var roleID uint
	_, err := fmt.Sscanf(id, "%d", &roleID)
	if err != nil {
		return nil, errors.New("无效的角色ID")
	}
	role, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}
	vo := RoleToV0(role)
	return vo, nil
}

// UpdateRole 更新角色信息
func (s *RoleService) UpdateRole(role *models.Role) error {
	// 唯一性校验
	roles, _ := s.repo.ListRoles(role.Name, "", "", 1, 1)
	for _, r := range roles {
		if (r.Name == role.Name || r.Code == role.Code) && r.ID != role.ID {
			return fmt.Errorf("角色名称 '%s' 或者角色编码 '%s' 已存在", role.Name, role.Code)
		}
	}
	return s.repo.UpdateRole(role)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id string) error {
	var roleID uint
	_, err := fmt.Sscanf(id, "%d", &roleID)
	if err != nil {
		return errors.New("无效的角色ID")
	}
	// 检查是否有关联用户/权限等，可扩展repo方法
	return s.repo.DeleteRole(roleID)
}

// GetRoleMenus 获取角色菜单
func (s *RoleService) GetRoleMenus(roleID string) ([]uint, error) {
	var rid uint
	_, err := fmt.Sscanf(roleID, "%d", &rid)
	if err != nil {
		return nil, errors.New("无效的角色ID")
	}
	return s.repo.GetRoleMenus(rid)
}

// UpdateRoleMenu 更新角色菜单
func (s *RoleService) UpdateRoleMenu(roleID uint, menuIDs []uint) error {
	return s.repo.UpdateRoleMenu(roleID, menuIDs)
}

func (s *RoleService) GetRoleOptions() ([]models.OptionLong, error) {
	roles, err := s.repo.ListRoles("", "", "", 1, 1000)
	if err != nil {
		return nil, err
	}
	return buildRoleOptions(roles), nil
}

func buildRoleOptions(roles []models.Role) []models.OptionLong {
	var options []models.OptionLong
	for _, role := range roles {
		option := models.OptionLong{
			Label: role.Name,
			Value: role.ID,
		}
		options = append(options, option)
	}
	return options
}

// RoleToV0 单个转换
func RoleToV0(r *models.Role) models.RoleV0 {
	return models.RoleV0{
		ID:          r.ID,
		Name:        r.Name,
		Code:        r.Code,
		Status:      r.Status,
		Sort:        r.Sort,
		DataScope:   r.DataScope,
		Description: "", // 如有描述字段可补充
	}
}

// RolesToV0List 批量转换
func RolesToV0List(roles []models.Role) []models.RoleV0 {
	result := make([]models.RoleV0, 0, len(roles))
	for _, r := range roles {
		result = append(result, RoleToV0(&r))
	}
	return result
}

// GetRolePerms 获取角色权限 code 列表
func (s *RoleService) GetRolePerms(roleID string) ([]string, error) {
	var rid uint
	_, err := fmt.Sscanf(roleID, "%d", &rid)
	if err != nil {
		return nil, errors.New("无效的角色ID")
	}
	return s.repo.GetRolePermCodes(rid)
}

// UpdateRolePerms 更新角色权限（code数组）
func (s *RoleService) UpdateRolePerms(roleID uint, permCodes []string) error {
	return s.repo.UpdateRolePerms(roleID, permCodes)
}
