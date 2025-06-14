package repositories

import (
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) UpdateRole(role *models.Role) error {
	updateMap := map[string]interface{}{
		"name":       role.Name,
		"code":       role.Code,
		"sort":       role.Sort,
		"status":     role.Status,
		"data_scope": role.DataScope,
	}
	return r.db.Model(&models.Role{}).Where("id = ?", role.ID).Updates(updateMap).Error
}

func (r *RoleRepository) DeleteRole(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *RoleRepository) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) ListRoles(keywords, startDate, endDate string, pageNum, pageSize int) ([]models.Role, error) {
	var roles []models.Role
	query := r.db.Model(&models.Role{})
	if keywords != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}
	if err := query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("sort ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) GetRoleMenus(roleID uint) ([]uint, error) {
	var roleMenus []models.RoleMenu
	if err := r.db.Table("role_menu").Where("role_id = ?", roleID).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	menuIDs := make([]uint, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}
	return menuIDs, nil
}

func (r *RoleRepository) UpdateRoleMenu(roleID uint, menuIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("role_menu").Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		var roleMenus []models.RoleMenu
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, models.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			})
		}
		if len(roleMenus) > 0 {
			if err := tx.Table("role_menu").Create(&roleMenus).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetRolePermCodes 查询角色权限 code 列表（只查role_permissions表）
func (r *RoleRepository) GetRolePermCodes(roleID uint) ([]string, error) {
	var codes []string
	err := r.db.Table("role_permissions").
		Where("role_id = ?", roleID).
		Pluck("permission_code", &codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}

// UpdateRolePerms 更新角色权限（permission_code数组）
func (r *RoleRepository) UpdateRolePerms(roleID uint, permCodes []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("role_permissions").Where("role_id = ?", roleID).Delete(nil).Error; err != nil {
			return err
		}
		var rolePerms []map[string]interface{}
		for _, code := range permCodes {
			rolePerms = append(rolePerms, map[string]interface{}{
				"role_id":         roleID,
				"permission_code": code,
			})
		}
		if len(rolePerms) > 0 {
			if err := tx.Table("role_permissions").Create(&rolePerms).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 角色菜单相关方法、角色选项等可后续补充
