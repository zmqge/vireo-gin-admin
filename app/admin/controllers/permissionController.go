package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

// @Group(path="/api/v1", name="权限管理", desc="权限相关接口")
type PermissionController struct {
	permissionService *services.PermissionService
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: services.NewPermissionService(),
	}
}

// 获取权限树
func (c *PermissionController) Tree(ctx *gin.Context) {
	tree, err := c.permissionService.GetPermissionTree()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, tree)
}

// 同步权限（从代码扫描）
func (c *PermissionController) Sync(ctx *gin.Context) {
	if err := c.permissionService.SyncPermissions(); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// ListPermOptions 获取权限下拉列表
// @Route(method=GET, path="/perms/options", middlewares=["jwt"])
// @Permission(code="sys:perm:options", name="权限下拉列表", modules="角色管理", desc="获取权限下拉选项")
func (c *PermissionController) ListPermOptions(ctx *gin.Context) {
	options, err := c.permissionService.ListPermOptions()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, options)
}
