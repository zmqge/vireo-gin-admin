package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
	"github.com/zmqge/vireo-gin-admin/utils"
)

// NoticesController Notices控制器
// @Group(path="/api/v1/", name="Notices管理")
type NoticesController struct {
	service services.NoticesService
}

// NewNoticesController 创建Notices控制器
func NewNoticesController(service services.NoticesService) *NoticesController {
	return &NoticesController{service: service}
}

// getNotices 获取单个Notices
// @Route(method=GET, path="/notices/:id/detail", middlewares=["jwt","dataperm"])
// @Permission(code="sys:notice:detail", name="Notices详情",modules="Notices管理", desc="查看Notices详情")
func (c *NoticesController) GetNoticesDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	entity, err := c.service.GetNoticesByID(ctx, uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// listNoticess 获取Notices分页列表
// @Route(method=GET, path="/notices/page", middlewares=["jwt","dataperm"])
// @Permission(code="sys:notice:query",name="Notices列表",modules="Notices管理", desc="查看Notices列表")
func (c *NoticesController) ListNoticess(ctx *gin.Context) {
	keywords := ctx.Query("keywords")
	publishStatus := ctx.Query("publishStatus")
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

	list, total, err := c.service.PageNotices(ctx, keywords, publishStatus, pageNum, pageSize)
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

// createNotices 创建Notices
// @Route(method=POST, path="/notices", middlewares=["jwt"])
// @Permission(code="sys:notice:add",name="新建Notices",modules="Notices管理", desc="创建Notices")
func (c *NoticesController) CreateNotices(ctx *gin.Context) {
	var entity models.NoticesModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	err := c.service.CreateNotices(ctx, &entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// updateNotices 更新Notices
// @Route(method=PUT, path="/notices/:id", middlewares=["jwt","dataperm"])
// @Permission(code="sys:notice:update",name="更新Notices",modules="Notices管理", desc="更新Notices")
func (c *NoticesController) UpdateNotices(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	var entity models.NoticesModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	entity.ID = uint(id)
	err = c.service.UpdateNotices(ctx, &entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// deleteNotices 删除Notices
// @Route(method=DELETE, path="/notices/:id", middlewares=["jwt"])
// @Permission(code="sys:notice:delete",name="删除Notices",modules="Notices管理", desc="删除Notices")
func (c *NoticesController) DeleteNotices(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.DeleteNotices(ctx, uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// getDict 获取Notices表单
// @Route(method=GET, path="/notices/:id/form", middlewares=["jwt","dataperm"])
// @Permission(code="sys:notice:form",name="Notices表单",modules="Notices管理", desc="查看Notices详情")
func (c *NoticesController) GetNoticesForm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}
	entity, err := c.service.GetNoticesByID(ctx, uint(id))

	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// revokeNotice 撤销公告
// @Route(method=PUT, path="/notices/:id/revoke", middlewares=["jwt"])
// @Permission(code="sys:notice:revoke",name="撤销公告",modules="Notices管理", desc="撤销公告")
func (c *NoticesController) RevokeNotice(ctx *gin.Context) {
	id, err := utils.ParseUintID(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.RevokeNotice(ctx, id)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// publishNotice 发布公告
// @Route(method=PUT, path="/notices/:id/publish", middlewares=["jwt"])
// @Permission(code="sys:notice:publish",name="发布公告",modules="Notices管理", desc="发布公告")
func (c *NoticesController) PublishNotice(ctx *gin.Context) {
	id, err := utils.ParseUintID(ctx.Param("id"))
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.PublishNoticeWithReceivers(ctx, id)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// 获取我的公告列表
// @Route(method=GET, path="/notices/my-page", middlewares=["jwt","dataperm"])
// @Permission(code="sys:notice:mynotice",name="我的公告列表",modules="Notices管理", desc="查看我的列表")
func (c *NoticesController) GetMyNoticess(ctx *gin.Context) {
	keywords := ctx.Query("keywords")
	publishStatus := ctx.Query("publishStatus")
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

	list, total, err := c.service.PageNotices(ctx, keywords, publishStatus, pageNum, pageSize)
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
