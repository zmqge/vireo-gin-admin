package repositories

import (
	"errors"
	"fmt"

	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetMenuByID(id string) (*models.Menu, error) {
	var menu models.Menu
	if err := r.db.Table("menu").Where("id = ?", id).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) ListMenus(keywords string) ([]models.Menu, error) {
	var menus []models.Menu
	query := r.db.Table("menu")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%").
			Or("name LIKE ?", "%"+keywords+"%")
	}
	if err := query.Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *MenuRepository) CreateMenu(menu *models.Menu) error {
	return r.db.Table("menu").Create(menu).Error
}

func (r *MenuRepository) UpdateMenu(menu *models.Menu) error {
	updateMap := map[string]interface{}{
		"parent_id":   menu.ParentID,
		"title":       menu.Title,
		"type":        menu.Type,
		"name":        menu.Name,
		"path":        menu.Path,
		"component":   menu.Component,
		"perm":        menu.Perm,
		"params":      menu.Params,
		"visible":     menu.Visible,
		"sort":        menu.Sort,
		"icon":        menu.Icon,
		"redirect":    menu.Redirect,
		"keep_alive":  menu.KeepAlive,
		"always_show": menu.AlwaysShow,
	}
	return r.db.Table("menu").Where("id =?", menu.ID).Updates(updateMap).Error
}

func (r *MenuRepository) DeleteMenu(id string) error {
	var menu models.Menu
	if err := r.db.Table("menu").Where("id = ?", id).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("菜单不存在或已被删除")
		}
		return fmt.Errorf("查询菜单失败: %v", err)
	}
	var childCount int64
	if err := r.db.Table("menu").Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return fmt.Errorf("检查子菜单失败: %v", err)
	}
	if childCount > 0 {
		return fmt.Errorf("无法删除菜单，仍有 %d 个子菜单存在", childCount)
	}
	var roleCount int64
	if err := r.db.Table("role_menu").Where("menu_id = ?", id).Count(&roleCount).Error; err != nil {
		return fmt.Errorf("检查菜单关联角色失败: %v", err)
	}
	if roleCount > 0 {
		return fmt.Errorf("无法删除菜单，仍有 %d 个角色关联此菜单", roleCount)
	}
	result := r.db.Table("menu").Where("id = ?", id).Delete(&models.Menu{})
	if result.Error != nil {
		return fmt.Errorf("删除菜单失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("菜单不存在或已被删除")
	}
	return nil
}

func (r *MenuRepository) ListMenuOptions() ([]models.Menu, error) {
	var menus []models.Menu
	if err := r.db.Table("menu").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *MenuRepository) ListUserRoleIDs(userID string) ([]uint, error) {
	var roleIDs []uint
	if err := r.db.Table("user_roles").Where("user_id = ?", userID).Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	return roleIDs, nil
}

func (r *MenuRepository) ListMenusByRoleIDs(roleIDs []uint) ([]models.Menu, error) {
	var menus []models.Menu
	if len(roleIDs) == 0 {
		return menus, nil
	}
	if err := r.db.Table("menu").
		Select("DISTINCT menu.*").
		Joins("LEFT JOIN role_menu ON role_menu.menu_id = menu.id").
		Where("role_menu.role_id IN ?", roleIDs).
		Where("menu.type = ?", 1).
		Order("menu.sort ASC").
		Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
