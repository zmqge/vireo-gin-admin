package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

// ConfigController Config控制器
// @Group(path="/api/v1/", name="Config管理")
type ConfigController struct {
	service services.ConfigService
}

// NewConfigController 创建Config控制器
func NewConfigController(service services.ConfigService) *ConfigController {
	return &ConfigController{service: service}
}

// getConfig 获取单个Config
// @Route(method=GET, path="/config/:id", middlewares=["jwt","dataperm"])
// @Permission(code="sys:config:view", name="Config详情",modules="Config管理", desc="查看Config详情")
func (c *ConfigController) GetConfigDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	entity, err := c.service.GetConfigByID(ctx, uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// listConfigs 获取Config分页列表
// @Route(method=GET, path="/config/page", middlewares=["jwt","dataperm"])
// @Permission(code="sys:config:query",name="Config列表",modules="Config管理", desc="查看Config列表")
func (c *ConfigController) ListConfigs(ctx *gin.Context) {
	keywords := ctx.Query("keywords")
	pageNumStr := ctx.DefaultQuery("pageNum", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")
	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	list, total, err := c.service.PageConfigs(ctx, keywords, pageNum, pageSize)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	resp := map[string]interface{}{
		"list":  list,
		"total": total,
	}
	response.Success(ctx, resp)
}

// createConfig 创建Config
// @Route(method=POST, path="/config", middlewares=["jwt"])
// @Permission(code="sys:config:add",name="新建Config",modules="Config管理", desc="创建Config")
func (c *ConfigController) CreateConfig(ctx *gin.Context) {
	var entity models.ConfigModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	err := c.service.CreateConfig(ctx, &entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// updateConfig 更新Config
// @Route(method=PUT, path="/config/:id", middlewares=["jwt","dataperm"])
// @Permission(code="sys:config:update",name="更新Config",modules="Config管理", desc="更新Config")
func (c *ConfigController) UpdateConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	var entity models.ConfigModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	entity.ID = uint(id)
	err = c.service.UpdateConfig(ctx, &entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// deleteConfig 删除Config
// @Route(method=DELETE, path="/config/:id", middlewares=["jwt"])
// @Permission(code="sys:config:delete",name="删除Config",modules="Config管理", desc="删除Config")
func (c *ConfigController) DeleteConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.DeleteConfig(ctx, uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// getDict 获取Config表单
// @Route(method=GET, path="/config/:id/form", middlewares=["jwt","dataperm"])
// @Permission(code="sys:config:details",name="Config详情",modules="Config管理", desc="查看Config详情")
func (c *ConfigController) GetConfigForm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}
	entity, err := c.service.GetConfigByID(ctx, uint(id))

	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}
