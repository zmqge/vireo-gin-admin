package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
	"github.com/zmqge/vireo-gin-admin/utils"
	"gorm.io/gorm"
)

// @Group(path="/api/v1/", name="部门管理")
type DeptController struct {
	deptService *services.DeptService
}

// NewDeptController 创建一个新的部门控制器实例。
func NewDeptController(db *gorm.DB) *DeptController {
	return &DeptController{deptService: services.NewDeptService(db)}
}

// 创建部门
// @Summary 创建部门
// @Route(method="POST", path="/dept",  middlewares=["jwt"])
// @Permission(code="sys:dept:add", name="创建部门", modules="部门管理", desc="创建部门")
func (c *DeptController) CreateDept(ctx *gin.Context) {
	var form struct {
		Name     string `json:"name" binding:"required"` // 部门名称
		ParentID uint   `json:"parentId"`                // 父部门ID，默认为0（根部门）
		Status   int    `json:"status"`                  // 状态，默认为1（启用）
		Sort     int    `json:"sort"`                    // 排序，默认为0
		Code     string `json:"code"`                    // 部门编码
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.Error(ctx, errors.New("invalid request body: "+err.Error()))
		return
	}

	// 调用服务层创建部门
	dept := &models.Dept{
		Name:     form.Name,
		ParentID: form.ParentID,
		Status:   form.Status,
		Sort:     form.Sort,
		Code:     form.Code,
	}
	if err := c.deptService.CreateDept(dept); err != nil {
		logrus.Errorf("Failed to create dept: %v", err)
		response.Error(ctx, errors.New("failed to create dept"))
		return
	}

	response.Success(ctx, gin.H{
		"id": dept.ID,
	})
}

// 更新部门
// @Route(method="PUT", path="/dept/:id",  middlewares=["jwt"])
// @Permission(code="sys:dept:edit", name="编辑部门", modules="部门管理", desc="更新部门")
func (c *DeptController) UpdateDept(ctx *gin.Context) {
	var form struct {
		Name     string `json:"name"`     // 部门名称
		ParentID uint   `json:"parentId"` // 父部门ID，默认为0（根部门）
		Status   int    `json:"status"`   // 状态，默认为1（启用）
		Sort     int    `json:"sort"`     // 排序，默认为0
		Code     string `json:"code"`     // 部门编码
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.Error(ctx, errors.New("invalid request body: "+err.Error()))
		return
	}

	// 获取部门ID
	id := ctx.Param("id")
	// 由于 id 是字符串类型，需要将其转换为 uint 类型
	// 定义 deptID 变量并使用 utils 包的函数将字符串 id 转换为 uint 类型
	var deptID uint
	var err error
	if deptID, err = utils.ParseUintID(id); err != nil {
		logrus.Errorf("Failed to convert id to uint: %v", err)
		response.Error(ctx, errors.New("invalid id"))
		return
	}

	// 调用服务层更新部门
	dept := &models.Dept{
		ID:       deptID, // 使用路径参数中的ID
		Name:     form.Name,
		ParentID: form.ParentID,
		Status:   form.Status,
		Sort:     form.Sort,
		Code:     form.Code,
	}
	if err := c.deptService.UpdateDept(ctx, dept); err != nil {
		logrus.Errorf("Failed to update dept: %v", err)
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// 获取部门详情
// @Route(method="GET", path="/dept/:id/form",  middlewares=["jwt"])
// @Permission(code="sys:dept:view", name="查看部门", modules="部门管理", desc="获取部门详情")
func (c *DeptController) GetDept(ctx *gin.Context) {
	// 获取部门ID
	id := ctx.Param("id")
	// 由于 id 是字符串类型，需要将其转换为 uint 类型
	var deptID uint
	var err error
	if deptID, err = utils.ParseUintID(id); err != nil {
		logrus.Errorf("Failed to convert id to uint: %v", err)
	}

	// 调用服务层获取部门详情
	dept, err := c.deptService.GetDeptDetails(deptID)
	if err != nil {
		logrus.Errorf("Failed to get dept: %v", err)
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, dept)
}

// 获取部门下拉列表
// @Route(method="GET", path="/dept/options",  middlewares=["jwt","dataperm"])
// @Permission(code="sys:dept:options", name="部门下拉列表", modules="部门管理", desc="获取部门下拉列表")
func (c *DeptController) GetDeptOptions(ctx *gin.Context) {
	// 调用服务层获取菜单列表
	menus, err := c.deptService.GetDeptOptions(ctx)
	if err != nil {
		logrus.Errorf("Failed to get depts: %v", err)
		response.Error(ctx, errors.New("failed to fetch depts"))
		return
	}
	response.Success(ctx, menus)
}

// @Summary 获取部门列表
// @Description 获取部门列表
// @Tags 部门管理
// @Route(method="GET", path="/dept",  middlewares=["jwt"])
// @Permission(code="sys:dept:query", name="部门列表", modules="部门管理", desc="获取部门列表")
func (c *DeptController) GetDepts(ctx *gin.Context) {
	// 获取查询参数
	keywords := ctx.Query("keywords")
	status := ctx.Query("status")

	// 调用服务层获取菜单列表
	menus, err := c.deptService.GetDepts(ctx, keywords, status)
	if err != nil {
		logrus.Errorf("Failed to get depts: %v", err)
		response.Error(ctx, errors.New("failed to fetch depts"))
		return
	}

	response.Success(ctx, menus)
}

// 删除部门
// @Route(method=DELETE, path="/dept/:id",  middlewares=["jwt"])
// @Permission(code="sys:dept:delete", name="删除部门", modules="部门管理", desc="删除部门")
func (c *DeptController) DeleteDept(ctx *gin.Context) {
	// 获取部门ID
	id := ctx.Param("id")
	// 由于 id 是字符串类型，需要将其转换为 uint 类型
	var deptID uint
	var err error
	if deptID, err = utils.ParseUintID(id); err != nil {
		logrus.Errorf("Failed to convert id to uint: %v", err)
		response.Error(ctx, errors.New("invalid id"))
		return
	}

	// 调用服务层删除部门
	if err := c.deptService.DeleteDept(ctx, deptID); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}
