package services

import (
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
	"gorm.io/gorm"
)

type MenuService struct {
	db   *gorm.DB
	repo *repositories.MenuRepository
}

func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{db: db, repo: repositories.NewMenuRepository(db)}
}

// GetCurrentUserRoutes 获取当前用户的路由列表（分级结构）
func (s *MenuService) GetCurrentUserRoutes(userID string) ([]models.RouteVO, error) {
	roleIDs, err := s.repo.ListUserRoleIDs(userID)
	if err != nil {
		return nil, err
	}
	menus, err := s.repo.ListMenusByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}
	return s.buildRouteTree(menus, 0), nil
}

// buildRouteTree 构建路由树形结构
func (s *MenuService) buildRouteTree(menus []models.Menu, parentID uint) []models.RouteVO {
	var routes []models.RouteVO

	for _, menu := range menus {
		if menu.ParentID == parentID {
			route := models.RouteVO{
				Path:      menu.Path,
				Component: menu.Component,
				Redirect:  menu.Redirect,
				Name:      menu.Name,
				ParentID:  menu.ParentID,
				// 修复 Meta 未定义的问题，推测 Meta 是 models 包中的类型，添加 models 包引用前缀
				Meta: models.Meta{
					Title: menu.Title,
					Icon:  menu.Icon,
					// 由于 menu.Visible 是 int 类型，不能直接使用 ! 操作符，这里假设值为 0 表示隐藏，非 0 表示显示
					Hidden:     menu.Visible == 0,
					KeepAlive:  menu.KeepAlive == 1,
					AlwaysShow: menu.AlwaysShow == 1,
				},
			}

			// 递归查找子路由
			children := s.buildRouteTree(menus, menu.ID)
			if len(children) > 0 {
				route.Children = &children

				// 如果父路由没有组件，设置为Layout
				if route.Component == "" {
					route.Component = "Layout"
				}
			}
			// 没有子路由时Children保持nil

			routes = append(routes, route)
		}
	}

	return routes
}

// GetMenus 获取菜单列表
func (s *MenuService) GetMenus(keywords string) ([]models.MenuVO, error) {
	menus, err := s.repo.ListMenus(keywords)
	if err != nil {
		return nil, err
	}
	menuTree := s.GetDynamicMenuTree(menus)
	return menuTree, nil
}

// buildMenuTree 构建菜单树形结构
func (s *MenuService) buildMenuTree(menus []models.Menu, parentID uint) []models.MenuVO {
	var result []models.MenuVO
	// 构建菜单ID索引（方便快速查找子节点）
	menuMap := make(map[uint]models.Menu)
	for _, m := range menus {
		menuMap[m.ID] = m
	}

	// 遍历所有菜单，找到父节点为当前parentID的节点
	for _, menu := range menus {
		if menu.ParentID == parentID {
			// 转换为VO对象
			menuVO := models.MenuVO{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Name:      menu.Title,
				Type:      menu.Type,
				RouteName: menu.Name,
				RoutePath: menu.Path,
				Component: menu.Component,
				Sort:      menu.Sort,
				Visible:   menu.Visible,
				Icon:      menu.Icon,
				Redirect:  menu.Redirect,
			}

			// 递归查找子节点（子节点的ParentID等于当前菜单的ID）
			children := s.buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menuVO.Children = &children
			}

			result = append(result, menuVO)
		}
	}

	return result
}

// 新增入口函数：动态获取所有根节点并构建树
func (s *MenuService) GetDynamicMenuTree(menus []models.Menu) []models.MenuVO {
	// 提取所有菜单的ID列表
	allMenuIDs := make(map[uint]bool)
	for _, m := range menus {
		allMenuIDs[m.ID] = true
	}

	// 查找所有根节点：ParentID不存在于allMenuIDs中，或ParentID为0（兼容旧逻辑）
	var rootParentIDs []uint
	for _, m := range menus {
		parentID := m.ParentID
		// 判断parentID是否是有效节点（即是否存在于allMenuIDs中）
		if !allMenuIDs[parentID] && parentID != 0 {
			// 自定义根节点条件：若parentID不是有效节点，则视为根节点
			rootParentIDs = append(rootParentIDs, parentID)
		}
	}

	// 去重根节点的ParentID（避免重复处理）
	rootSet := make(map[uint]bool)
	for _, pid := range rootParentIDs {
		rootSet[pid] = true
	}
	var uniqueRoots []uint
	for pid := range rootSet {
		uniqueRoots = append(uniqueRoots, pid)
	}

	// 构建菜单树：若存在自定义根节点，优先处理；否则默认从parentID=0开始
	var result []models.MenuVO
	if len(uniqueRoots) > 0 {
		for _, pid := range uniqueRoots {
			result = append(result, s.buildMenuTree(menus, pid)...)
		}
	} else {
		// 兼容旧逻辑：无自定义根节点时，从parentID=0开始
		result = s.buildMenuTree(menus, 0)
	}

	return result
}

func buildMenuOptions(menus []models.Menu, parentID uint) []models.OptionLong {
	var options []models.OptionLong

	for _, menu := range menus {
		if menu.ParentID == parentID {
			option := models.OptionLong{
				Value: menu.ID,
				Label: menu.Title,
			}

			children := buildMenuOptions(menus, menu.ID)
			if len(children) > 0 {
				option.Children = children
			}

			options = append(options, option)
		}
	}

	return options
}

func (s *MenuService) GetMenuOptions(onlyParent bool) ([]models.OptionLong, error) {
	menus, err := s.repo.ListMenuOptions()
	if err != nil {
		return nil, err
	}
	return buildMenuOptions(menus, 0), nil
}

func (s *MenuService) GetMenuDetail(id string) (*models.Menu, error) {
	return s.repo.GetMenuByID(id)
}

// CreateMenu 创建新菜单
func (s *MenuService) CreateMenu(menu *models.Menu) error {
	return s.repo.CreateMenu(menu)
}

// UpdateMenu 更新菜单信息
func (s *MenuService) UpdateMenu(menu *models.Menu) error {
	return s.repo.UpdateMenu(menu)
}

// DeleteMenu 删除菜单
// DeleteMenu 删除菜单（带关联检查的硬删除）
func (s *MenuService) DeleteMenu(id string) error {
	return s.repo.DeleteMenu(id)
}
