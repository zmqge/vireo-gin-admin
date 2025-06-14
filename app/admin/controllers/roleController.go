package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
	"gorm.io/gorm"
)

// @Group(path="/api/v1/", name="角色管理")
type RoleController struct {
	roleService *services.RoleService
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{roleService: services.NewRoleService(db)}

}

// 创建角色
// @Route(method="POST", path="/roles", middlewares=["jwt"])
// @Permission(code="sys:role:add",name="新增角色",modules="角色管理", desc="创建角色")
func (c *RoleController) Create(ctx *gin.Context) {
	var input struct {
		Name        string   `json:"name" binding:"required"`
		Code        string   `json:"code" binding:"required"`
		Sort        int      `json:"sort"`
		Status      int      `json:"status"`
		DataScope   int      `json:"dataScope"`
		Permissions []string `json:"permissions"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	role := models.Role{
		Name:      input.Name,
		Code:      input.Code,
		Sort:      input.Sort,
		Status:    input.Status,
		DataScope: input.DataScope,
	}
	err := c.roleService.CreateRole(&role)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// 更新角色
// @Route(method="PUT", path="/roles/:id", middlewares=["jwt"])
// @Permission(code="sys:role:edit",name="编辑角色",modules="角色管理", desc="更新角色")
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	var input struct {
		Name        string   `json:"name" binding:"required"`
		Code        string   `json:"code" binding:"required"`
		Sort        int      `json:"sort"`
		Status      int      `json:"status"`
		DataScope   int      `json:"dataScope"`
		Permissions []string `json:"permissions"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	role := models.Role{
		Name:      input.Name,
		Code:      input.Code,
		Sort:      input.Sort,
		Status:    input.Status,
		DataScope: input.DataScope,
	}
	// 由于 c.roleService.UpdateRole 方法只接受 *models.Role 类型参数，移除 ctx.Param("id") 参数
	// 假设角色的 ID 已经在 role 结构体中设置，直接传递 role 指针
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, "无效的角色 ID")
		return
	}
	role.ID = uint(id)
	err = c.roleService.UpdateRole(&role)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// 删除角色
// @Route(method="DELETE", path="/roles/:id", middlewares=["jwt"])
// @Permission(code="sys:role:delete",name="删除角色",modules="角色管理", desc="删除角色")
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	err := c.roleService.DeleteRole(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// 获取角色详情
// @Route(method="GET", path="/roles/:id/form", middlewares=["jwt"])
// @Permission(code="sys:role:detail",name="角色详情",modules="角色管理", desc="获取角色详情")
func (c *RoleController) GetRoleDetail(ctx *gin.Context) {
	role, err := c.roleService.GetRoleDetails(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, role)
}

// 获取角色分页列表
// @Route(method="GET", path="/roles/page", middlewares=["jwt"])
// @Permission(code="sys:role:query",name="角色列表",modules="角色管理", desc="获取角色分页列表")
func (c *RoleController) List(ctx *gin.Context) {
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	keywords := ctx.Query("keywords")
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	// 由于 c.roleService.GetRolePage 只返回 2 个值，这里调整变量接收数量
	roles, err := c.roleService.GetRolePage(keywords, startDate, endDate, pageNum, pageSize)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, roles)

}

// GetRoleMenus 获取角色菜单
// @Route(method="GET", path="/roles/:id/menuIds", middlewares=["jwt"])
// @Permission(code="sys:role:menu",name="角色菜单列表",modules="角色管理", desc="获取角色菜单")
func (c *RoleController) GetRoleMenus(ctx *gin.Context) {
	roleID := ctx.Param("id")
	if roleID == "" {
		response.BadRequest(ctx, "无效的角色 ID")
		return
	}
	// 由于 c.roleService.GetRoleMenus 可能需要字符串类型参数，将 uint 类型的 roleID 转换为字符串
	menuIDs, err := c.roleService.GetRoleMenus(roleID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, menuIDs)
}

// GetRolePerms 获取角色权限 code 列表
// @Route(method="GET", path="/roles/:id/permCodes", middlewares=["jwt"])
// @Permission(code="sys:role:perm",name="角色权限列表",modules="角色管理", desc="获取角色权限")
func (c *RoleController) GetRolePerms(ctx *gin.Context) {
	roleID := ctx.Param("id")
	if roleID == "" {
		response.BadRequest(ctx, "无效的角色 ID")
		return
	}
	permCodes, err := c.roleService.GetRolePerms(roleID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, permCodes)
}

// UpdateRoleMenus 更新角色菜单
// @Route(method="PUT", path="/roles/:id/menus", middlewares=["jwt"])
// @Permission(code="sys:role:menu:update",name="更新角色菜单",modules="角色管理", desc="更新角色菜单")
func (c *RoleController) UpdateRoleMenus(ctx *gin.Context) {
	var input struct {
		MenuIDs []uint `json:"menuIds"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, "无效的角色 ID")
		return
	}
	// 由于 c.roleService.UpdateRoleMenus 可能需要字符串类型参数，将 uint 类型的 roleID 转换为字符串
	// 由于 UpdateRoleMenus 方法需要 uint 类型参数，直接使用 roleID 转换为 uint 类型
	err = c.roleService.UpdateRoleMenu(uint(roleID), input.MenuIDs)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// UpdateRolePerms 更新角色权限
// @Route(method="PUT", path="/roles/:id/perms", middlewares=["jwt"])
// @Permission(code="sys:role:perm:update",name="更新角色权限",modules="角色管理", desc="更新角色权限")
func (c *RoleController) UpdateRolePerms(ctx *gin.Context) {
	var input struct {
		PermCodes []string `json:"permCodes"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, "无效的角色 ID")
		return
	}
	// 过滤掉 "0"
	var codes []string
	for _, code := range input.PermCodes {
		if code != "0" {
			codes = append(codes, code)
		}
	}
	err = c.roleService.UpdateRolePerms(uint(roleID), codes)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// ListRoleOptions 获取角色下拉列表
// @Route(method=GET, path="/roles/options", middlewares=["jwt"])
// @Permission(code="sys:role:options",name="角色下拉列表",modules="角色管理", desc="获取角色下拉列表")
func (c *RoleController) ListRoleOptions(ctx *gin.Context) {

	// 调用服务层获取菜单下拉列表
	options, err := c.roleService.GetRoleOptions()
	if err != nil {
		logrus.Errorf("Failed to get role options: %v", err)
		response.Error(ctx, errors.New("failed to fetch role options"))
		return
	}

	response.Success(ctx, options)
}
