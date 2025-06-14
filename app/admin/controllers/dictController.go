package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

// DictController Dict控制器
// @Group(path="/api/v1/", desc="Dict相关接口")
type DictController struct {
	service services.DictService
}

// NewDictController 创建Dict控制器
func NewDictController(service services.DictService) *DictController {
	return &DictController{service: service}
}

// getDict 获取单个Dict
// @Route(method=GET, path="/dicts/:id/form", middlewares=["jwt"])
// @Permission(code="sys:dict:details",name="查看字典",modules="字典管理", desc="查看Dict详情")
func (c *DictController) GetDict(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	entity, err := c.service.GetDictByID(uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// listDicts 获取Dict分页列表
// @Route(method=GET, path="/dicts/page", middlewares=["jwt"])
// @Permission(code="sys:dict:query",name="字典查询",modules="字典管理", desc="查看Dict列表")
func (c *DictController) ListDicts(ctx *gin.Context) {
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

	list, total, err := c.service.PageDicts(keywords, pageNum, pageSize)
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

// createDict 创建Dict
// @Route(method=POST, path="/dicts", middlewares=["jwt"])
// @Permission(code="sys:dict:add",name="新增字典",modules="字典管理", desc="创建Dict")
func (c *DictController) CreateDict(ctx *gin.Context) {
	var entity models.DictModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	err := c.service.CreateDict(&entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// updateDict 更新Dict
// @Route(method=PUT, path="/dicts/:id", middlewares=["jwt"])
// @Permission(code="sys:dict:edit",name="编辑字典",modules="字典管理", desc="创建Dict")
func (c *DictController) UpdateDict(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	var entity models.DictModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	entity.ID = uint(id)
	err = c.service.UpdateDict(&entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// deleteDict 删除Dict
// @Route(method=DELETE, path="/dicts/:id", middlewares=["jwt"])
// @Permission(code="sys:dict-item:delete",name="删除字典项", desc="删除字典项",modules="字典管理")
func (c *DictController) DeleteDict(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.DeleteDict(uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// GetDictItem 获取Dict选项
// @Route(method=GET, path="/dicts-items/:dictCode/items", middlewares=["jwt"])
// @Permission(code="sys:dict-item:details",name="查看字典项",modules="字典项管理", desc="查看Dict选项")
func (c *DictController) GetDictItem(ctx *gin.Context) {
	dictCode := ctx.Param("dictCode")
	if dictCode == "" {
		response.BadRequest(ctx, "Dict code is required")
		return
	}

	items, err := c.service.GetDictItemsByCode(dictCode)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// OpenAPI 规范要求返回数组，且字段为 value/label/tagType
	var result []map[string]interface{}
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"value":   item.Value,
			"label":   item.Label,
			"tagType": item.TagType,
		})
	}
	response.Success(ctx, result)
}

// GetDictItemPage 获取字典项分页列表
// @Route(method=GET, path="/dicts-items/:dictCode/items/page", middlewares=["jwt"])
// @Permission(code="sys:dict-item:query",name="字典项查询",modules="字典项管理", desc="查看Dict选项")
func (c *DictController) GetDictItemPage(ctx *gin.Context) {
	dictCode := ctx.Param("dictCode")
	if dictCode == "" {
		response.BadRequest(ctx, "Dict code is required")
		return
	}
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

	list, total, err := c.service.PageDictItems(dictCode, keywords, pageNum, pageSize)
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

// CreateDictItem 新增字典项
// @Route(method=POST, path="/dicts-items/:dictCode/items", middlewares=["jwt"])
// @Permission(code="sys:dict-item:add",name="新增字典项",modules="字典项管理", desc="查看Dict选项")
func (c *DictController) CreateDictItem(ctx *gin.Context) {
	dictCode := ctx.Param("dictCode")
	if dictCode == "" {
		response.BadRequest(ctx, "Dict code is required")
		return
	}
	var form models.DictItemModel
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	form.DictCode = dictCode // 路径参数优先生效
	err := c.service.CreateDictItem(&form)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// UpdateDictItem 更新字典项
// @Route(method=PUT, path="/dicts-items/:dictCode/items/:id", middlewares=["jwt"])
// @Permission(code="sys:dict-item:edit",name="编辑字典项",modules="字典项管理", desc="查看Dict选项")
func (c *DictController) UpdateDictItem(ctx *gin.Context) {
	dictCode := ctx.Param("dictCode")
	if dictCode == "" {
		response.BadRequest(ctx, "Dict code is required")
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	var form models.DictItemModel
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	form.ID = uint(id)
	form.DictCode = dictCode // 路径参数优先生效

	err = c.service.UpdateDictItem(&form)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// DeleteDictItem 删除字典项
// @Route(method=DELETE, path="/dicts-items/:dictCode/items/:id", middlewares=["jwt"])
// @Permission(code="sys:dict:delete",name="删除字典项", desc="删除字典项",modules="字典项管理")
func (c *DictController) DeleteDictItem(ctx *gin.Context) {
	dictCode := ctx.Param("dictCode")
	if dictCode == "" {
		response.BadRequest(ctx, "Dict code is required")
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.DeleteDictItem(uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// GetDictItemForm 获取字典项表单数据
// @Route(method=GET, path="/dicts-items/:dictCode/items/:itemId/form", middlewares=["jwt"])
// @Permission(code="sys:dict-item:form",name="字典项表单",modules="字典项管理", desc="查看Dict选项")
func (c *DictController) GetDictItemForm(ctx *gin.Context) {
	dictCode := ctx.Param("dictCode")
	if dictCode == "" {
		response.BadRequest(ctx, "Dict code is required")
		return
	}
	itemIdStr := ctx.Param("itemId")
	itemId, err := strconv.ParseUint(itemIdStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid itemId")
		return
	}
	item, err := c.service.GetDictItemForm(dictCode, uint(itemId))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if item == nil {
		response.NotFound(ctx, "Dict item not found")
		return
	}

	response.Success(ctx, item)
}
