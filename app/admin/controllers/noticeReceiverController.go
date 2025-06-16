package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

// NoticeReceiverController NoticeReceiver控制器
// @Group(path="/api/v1/", name="NoticeReceiver管理")
type NoticeReceiverController struct {
	service services.NoticeReceiverService
}

// NewNoticeReceiverController 创建NoticeReceiver控制器
func NewNoticeReceiverController(service services.NoticeReceiverService) *NoticeReceiverController {
	return &NoticeReceiverController{service: service}
}

// getNoticeReceiver 获取单个NoticeReceiver
// @Route(method=GET, path="/noticereceiver/:id", middlewares=["jwt","dataperm"])
// @Permission(code="sys:noticereceiver:view", name="NoticeReceiver详情",modules="NoticeReceiver管理", desc="查看NoticeReceiver详情")
func (c *NoticeReceiverController) GetNoticeReceiverDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	entity, err := c.service.GetNoticeReceiverByID(ctx, uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// listNoticeReceivers 获取NoticeReceiver分页列表
// @Route(method=GET, path="/noticereceiver/page", middlewares=["jwt","dataperm"])
// @Permission(code="sys:noticereceiver:query",name="NoticeReceiver列表",modules="NoticeReceiver管理", desc="查看NoticeReceiver列表")
func (c *NoticeReceiverController) ListNoticeReceivers(ctx *gin.Context) {
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

	list, total, err := c.service.PageNoticeReceivers(ctx, keywords, pageNum, pageSize)
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

// createNoticeReceiver 创建NoticeReceiver
// @Route(method=POST, path="/noticereceiver", middlewares=["jwt"])
// @Permission(code="sys:noticereceiver:add",name="新建NoticeReceiver",modules="NoticeReceiver管理", desc="创建NoticeReceiver")
func (c *NoticeReceiverController) CreateNoticeReceiver(ctx *gin.Context) {
	var entity models.NoticeReceiverModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	err := c.service.CreateNoticeReceiver(ctx, &entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// updateNoticeReceiver 更新NoticeReceiver
// @Route(method=PUT, path="/noticereceiver/:id", middlewares=["jwt","dataperm"])
// @Permission(code="sys:noticereceiver:update",name="更新NoticeReceiver",modules="NoticeReceiver管理", desc="更新NoticeReceiver")
func (c *NoticeReceiverController) UpdateNoticeReceiver(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	var entity models.NoticeReceiverModel
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}
	entity.ID = uint(id)
	err = c.service.UpdateNoticeReceiver(ctx, &entity)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}

// deleteNoticeReceiver 删除NoticeReceiver
// @Route(method=DELETE, path="/noticereceiver/:id", middlewares=["jwt"])
// @Permission(code="sys:noticereceiver:delete",name="删除NoticeReceiver",modules="NoticeReceiver管理", desc="删除NoticeReceiver")
func (c *NoticeReceiverController) DeleteNoticeReceiver(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.DeleteNoticeReceiver(ctx, uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// getDict 获取NoticeReceiver表单
// @Route(method=GET, path="/noticereceiver/:id/form", middlewares=["jwt","dataperm"])
// @Permission(code="sys:noticereceiver:details",name="NoticeReceiver详情",modules="NoticeReceiver管理", desc="查看NoticeReceiver详情")
func (c *NoticeReceiverController) GetNoticeReceiverForm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid ID")
		return
	}
	entity, err := c.service.GetNoticeReceiverByID(ctx, uint(id))

	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, entity)
}
