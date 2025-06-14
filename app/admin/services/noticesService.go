package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
)

// NoticesService Notices服务接口
type NoticesService interface {
	GetNoticesByID(ctx *gin.Context, id uint) (*models.NoticesModel, error) // 使用uint类型
	ListNoticess(ctx *gin.Context) ([]*models.NoticesModel, error)
	CreateNotices(ctx *gin.Context, entity *models.NoticesModel) error
	UpdateNotices(ctx *gin.Context, entity *models.NoticesModel) error
	DeleteNotices(ctx *gin.Context, id uint) error // 使用uint类型
	PageNoticess(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error)
	GetNoticesForm(ctx *gin.Context, id uint) (*models.NoticesModel, error) // 使用uint类型
	RevokeNotice(ctx *gin.Context, id uint) error
	PublishNotice(ctx *gin.Context, id uint) error
}

// NoticesServiceImpl Notices服务实现
type NoticesServiceImpl struct {
	repo repositories.NoticesRepository
}

// NewNoticesService 创建Notices服务
func NewNoticesService(repo repositories.NoticesRepository) NoticesService {
	return &NoticesServiceImpl{repo: repo}
}

// GetNoticesByID 根据ID获取Notices
func (s *NoticesServiceImpl) GetNoticesByID(ctx *gin.Context, id uint) (*models.NoticesModel, error) {
	entity, err := s.repo.GetNoticesByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("Notices not found")
	}
	return entity, nil
}

// ListNoticess 获取Notices列表
func (s *NoticesServiceImpl) ListNoticess(ctx *gin.Context) ([]*models.NoticesModel, error) {
	return s.repo.ListNoticess(ctx)
}

// CreateNotices 创建Notices
func (s *NoticesServiceImpl) CreateNotices(ctx *gin.Context, entity *models.NoticesModel) error {
	if entity.Title == "" {
		return errors.New("name is required")
	}
	return s.repo.CreateNotices(ctx, entity)
}

// UpdateNotices 更新Notices
func (s *NoticesServiceImpl) UpdateNotices(ctx *gin.Context, entity *models.NoticesModel) error {
	if entity.Title == "" {
		return errors.New("name is required")
	}
	return s.repo.UpdateNotices(ctx, entity)
}

// DeleteNotices 删除Notices
func (s *NoticesServiceImpl) DeleteNotices(ctx *gin.Context, id uint) error {
	return s.repo.DeleteNotices(ctx, id)
}

// PageNoticess 分页获取Notices列表
func (s *NoticesServiceImpl) PageNoticess(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error) {
	return s.repo.PageNoticess(ctx, keywords, publishStatus, pageNum, pageSize)
}

// GetNoticesForm 表单
func (s *NoticesServiceImpl) GetNoticesForm(ctx *gin.Context, id uint) (*models.NoticesModel, error) {
	return s.repo.GetNoticesByID(ctx, id)
}

// RevokeNotice 撤回通知
func (s *NoticesServiceImpl) RevokeNotice(ctx *gin.Context, id uint) error {
	return s.repo.RevokeNotice(ctx, id)
}

// PublishNotice 发布通知
func (s *NoticesServiceImpl) PublishNotice(ctx *gin.Context, id uint) error {
	return s.repo.PublishNotice(ctx, id)
}
