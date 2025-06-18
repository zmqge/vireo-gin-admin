package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
	"gorm.io/gorm"
)

// @Group(name="菜单管理",path="/api/v1")
type MenuController struct {
	menuService *services.MenuService
}

func NewMenuController(db *gorm.DB) *MenuController {
	return &MenuController{
		menuService: services.NewMenuService(db),
	}
}

// GetCurrentUserRoutes 获取当前用户的路由列表
// @Route(method=GET, path="/menus/routes", middlewares=["jwt"])
func (c *MenuController) GetCurrentUserRoutes(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		response.Error(ctx, errors.New("userID is required"))
		return
	}
	routes, err := c.menuService.GetCurrentUserRoutes(userID)
	if err != nil {
		logrus.Errorf("Failed to get current user routes: %v", err)
		response.Error(ctx, errors.New("failed to fetch user routes"))
		return
	}

	response.Success(ctx, routes)
}

// ListMenus 获取菜单列表
// @Route(method=GET, path="/menus", middlewares=["jwt"])
// @Permission(code="sys:menu:query", name="菜单查询", modules="菜单管理", desc="查看菜单列表")
func (c *MenuController) ListMenus(ctx *gin.Context) {
	// 获取查询参数
	keywords := ctx.Query("keywords")

	// 调用服务层获取菜单列表
	menus, err := c.menuService.GetMenus(keywords)
	if err != nil {
		logrus.Errorf("Failed to get menus: %v", err)
		response.Error(ctx, errors.New("failed to fetch menus"))
		return
	}

	response.Success(ctx, menus)
}

// ListMenuOptions 获取菜单下拉列表
// @Route(method=GET, path="/menus/options", middlewares=["jwt"])
// @Permission(code="sys:menu:options", name="菜单下拉列表", modules="菜单管理", desc="获取菜单下拉选项")
func (c *MenuController) ListMenuOptions(ctx *gin.Context) {
	// 获取查询参数
	onlyParent := ctx.Query("onlyParent") == "true"

	// 调用服务层获取菜单下拉列表
	options, err := c.menuService.GetMenuOptions(onlyParent)
	if err != nil {
		logrus.Errorf("Failed to get menu options: %v", err)
		response.Error(ctx, errors.New("failed to fetch menu options"))
		return
	}

	response.Success(ctx, options)
}

// GetMenuDetail 获取菜单详情
// @Route(method=GET, path="/menus/:id/form", middlewares=["jwt"])
// @Permission(code="sys:menu:details", name="菜单详情", modules="菜单管理", desc="查看菜单详情")
func (c *MenuController) GetMenuDetail(ctx *gin.Context) {
	// 获取路径参数
	id := ctx.Param("id")
	if id == "" {
		response.Error(ctx, errors.New("menu id is required"))
		return
	}

	// 调用服务层获取菜单详情
	menu, err := c.menuService.GetMenuDetail(id)
	if err != nil {
		logrus.Errorf("Failed to get menu detail: %v", err)
		response.Error(ctx, errors.New("failed to fetch menu detail"))
		return
	}
	// jsonData, err := json.Marshal(menu.Params)
	// if err != nil {
	// 	logrus.Errorf("Failed to marshal menu params: %v", err)
	// 	jsonData = []byte("")
	// }

	var params []map[string]interface{}
	if menu.Params != "" {
		if err := json.Unmarshal([]byte(menu.Params), &params); err != nil {
			// 处理错误
		}
	} else {
		params = make([]map[string]interface{}, 0) // 初始化空数组
	}
	response.Success(ctx, map[string]interface{}{
		"id":         menu.ID,
		"parentId":   menu.ParentID,
		"name":       menu.Title,
		"type":       menu.Type,
		"routeName":  menu.Name,
		"routePath":  menu.Path,
		"component":  menu.Component,
		"perm":       menu.Perm,
		"visible":    menu.Visible,
		"sort":       menu.Sort,
		"icon":       menu.Icon,
		"redirect":   menu.Redirect,
		"keepAlive":  menu.KeepAlive,
		"alwaysShow": menu.AlwaysShow,
		"params":     params,

		// 数据库读取后的反序列化示例

	})
}

// AddMenu 新增菜单
// @Route(method=POST, path="/menus", middlewares=["jwt"])
// @Permission(code="sys:menu:add", name="新增菜单", modules="菜单管理", desc="新增菜单")
func (c *MenuController) AddMenu(ctx *gin.Context) {
	var form struct {
		ParentID   uint   `json:"parentId"`
		Name       string `json:"name"`
		Type       int    `json:"type"`
		RouteName  string `json:"routeName"`
		RoutePath  string `json:"routePath"`
		Component  string `json:"component"`
		Perm       string `json:"perm"`
		Visible    int    `json:"visible"`
		Sort       int    `json:"sort"`
		Icon       string `json:"icon"`
		Redirect   string `json:"redirect"`
		KeepAlive  int    `json:"keepAlive"`
		AlwaysShow int    `json:"alwaysShow"`
		Params     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"params"`
	}

	if err := ctx.ShouldBindJSON(&form); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal string into Go struct field .parentId of type uint") {
			// 尝试将字符串类型的parentId转换为uint
			// 从上下文中直接获取原始的 parentId 字符串值
			if parentIDStr, exists := ctx.GetPostForm("parentId"); exists {
				if parentID, _ := strconv.ParseUint(parentIDStr, 10, 64); err == nil {
					form.ParentID = uint(parentID)
				}
			}
			response.Error(ctx, errors.New("invalid parentId type"))
			return
		}
		logrus.Errorf("Invalid request body: %v", err)
		response.Error(ctx, errors.New("invalid request body: "+err.Error()))
		return
	}

	// Validate required fields
	if form.Name == "" || form.Type == 0 {
		response.Error(ctx, errors.New("name and type are required"))
		return
	}

	menu := models.Menu{
		ParentID:  form.ParentID,
		Title:     form.Name,
		Type:      form.Type,
		Name:      form.RouteName,
		Path:      form.RoutePath,
		Component: form.Component,
		Perm:      form.Perm,
		Params:    paramsToString(form.Params),
		Visible:   form.Visible,
		Sort:      form.Sort,
		Icon:      form.Icon,
		Redirect:  form.Redirect,
		// 将 int 类型的 form.KeepAlive 转换为 bool 类型，假设非零值为 true，零值为 false
		KeepAlive:  form.KeepAlive,
		AlwaysShow: form.AlwaysShow,
	}

	if err := c.menuService.CreateMenu(&menu); err != nil {
		logrus.Errorf("Failed to create menu: %v", err)
		response.Error(ctx, errors.New("failed to create menu"))
		return
	}

	response.Success(ctx, gin.H{
		"id": menu.ID,
	})
}

// 新增JSON序列化函数
func paramsToString(params []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}) string {
	if len(params) == 0 {
		return ""
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		logrus.WithField("params", params).Errorf("JSON序列化失败: %v", err)
		return ""
	}
	return string(jsonData)
}

// UpdateMenu 修改菜单
// @Route(method=PUT, path="/menus/:id", middlewares=["jwt"])
// @Permission(code="sys:menu:edit", name="编辑菜单", modules="菜单管理", desc="编辑菜单")
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	// 获取路径参数id
	// 由于变量 id 声明后未使用，暂时注释掉获取 id 的代码
	// id := ctx.Param("id")

	var form struct {
		ID         uint   `json:"id"`
		ParentID   uint   `json:"parentId"`
		Name       string `json:"name"`
		Type       int    `json:"type"`
		RouteName  string `json:"routeName"`
		RoutePath  string `json:"routePath"`
		Component  string `json:"component"`
		Perm       string `json:"perm"`
		Visible    int    `json:"visible"`
		Sort       int    `json:"sort"`
		Icon       string `json:"icon"`
		Redirect   string `json:"redirect"`
		KeepAlive  int    `json:"keepAlive"`
		AlwaysShow int    `json:"alwaysShow"`
		Params     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"params"`
	}

	if err := ctx.ShouldBindJSON(&form); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal string into Go struct field .parentId of type uint") {
			// 尝试将字符串类型的parentId转换为uint
			// 从上下文中直接获取原始的 parentId 字符串值
			if parentIDStr, exists := ctx.GetPostForm("parentId"); exists {
				if parentID, _ := strconv.ParseUint(parentIDStr, 10, 64); err == nil {
					form.ParentID = uint(parentID)
				}
			}
			response.Error(ctx, errors.New("invalid parentId type"))
			return
		}
		logrus.Errorf("Invalid request body: %v", err)
		response.Error(ctx, errors.New("invalid request body: "+err.Error()))
		return
	}

	// Validate required fields
	if form.Name == "" || form.Type == 0 {
		response.Error(ctx, errors.New("name and type are required"))
		return
	}

	menu := models.Menu{
		ID:        form.ID,
		ParentID:  form.ParentID,
		Title:     form.Name,
		Type:      form.Type,
		Name:      form.RouteName,
		Path:      form.RoutePath,
		Component: form.Component,
		Perm:      form.Perm,
		Params:    paramsToString(form.Params),
		Visible:   form.Visible,
		Sort:      form.Sort,
		Icon:      form.Icon,
		Redirect:  form.Redirect,
		// 将 int 类型的 form.KeepAlive 转换为 bool 类型，假设非零值为 true，零值为 false
		KeepAlive:  form.KeepAlive,
		AlwaysShow: form.AlwaysShow,
	}

	if err := c.menuService.UpdateMenu(&menu); err != nil {
		logrus.Errorf("Failed to update menu: %v", err)
		response.Error(ctx, errors.New("failed to update menu"))
		return
	}

	response.Success(ctx, gin.H{
		"id": menu.ID,
	})
}

// DeleteMenu 删除菜单
// @Route(method=DELETE, path="/menus/:id", middlewares=["jwt"])
// @Permission(code="sys:menu:delete", name="删除菜单", modules="菜单管理", desc="删除菜单")
func (c *MenuController) DeleteMenu(ctx *gin.Context) {

	// 调用服务层删除菜单
	err := c.menuService.DeleteMenu(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, nil)
}
