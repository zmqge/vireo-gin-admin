package services

import (
    "errors"
    "github.com/zmqge/vireo-gin-admin/app/admin/models"
    "github.com/zmqge/vireo-gin-admin/app/admin/repositories"
    "github.com/gin-gonic/gin"
)

// NoticeReceiverService NoticeReceiver服务接口
type NoticeReceiverService interface {
    GetNoticeReceiverByID(ctx *gin.Context,id uint) (*models.NoticeReceiverModel, error) // 使用uint类型
    ListNoticeReceivers(ctx *gin.Context,) ([]*models.NoticeReceiverModel, error)
    CreateNoticeReceiver(ctx *gin.Context,entity *models.NoticeReceiverModel) error
    UpdateNoticeReceiver(ctx *gin.Context,entity *models.NoticeReceiverModel) error
    DeleteNoticeReceiver(ctx *gin.Context,id uint) error // 使用uint类型
	PageNoticeReceivers(ctx *gin.Context,keywords string, pageNum, pageSize int) ([]*models.NoticeReceiverModel, int64, error)
	GetNoticeReceiverForm(ctx *gin.Context,id uint) (*models.NoticeReceiverModel, error)    // 使用uint类型
}

// NoticeReceiverServiceImpl NoticeReceiver服务实现
type NoticeReceiverServiceImpl struct {
    repo repositories.NoticeReceiverRepository
}

// NewNoticeReceiverService 创建NoticeReceiver服务
func NewNoticeReceiverService(repo repositories.NoticeReceiverRepository) NoticeReceiverService {
    return &NoticeReceiverServiceImpl{repo: repo}
}

// GetNoticeReceiverByID 根据ID获取NoticeReceiver
func (s *NoticeReceiverServiceImpl) GetNoticeReceiverByID(ctx *gin.Context,id uint) (*models.NoticeReceiverModel, error) {
    entity, err := s.repo.GetNoticeReceiverByID(ctx,id)
    if err != nil {
        return nil, err
    }
    if entity == nil {
        return nil, errors.New("NoticeReceiver not found")
    }
    return entity, nil
}

// ListNoticeReceivers 获取NoticeReceiver列表
func (s *NoticeReceiverServiceImpl) ListNoticeReceivers(ctx *gin.Context,) ([]*models.NoticeReceiverModel, error) {
    return s.repo.ListNoticeReceivers(ctx,)
}

// CreateNoticeReceiver 创建NoticeReceiver
func (s *NoticeReceiverServiceImpl) CreateNoticeReceiver(ctx *gin.Context,entity *models.NoticeReceiverModel) error {
    if entity.Name == "" {
        return errors.New("name is required")
    }
    return s.repo.CreateNoticeReceiver(ctx,entity)
}

// UpdateNoticeReceiver 更新NoticeReceiver
func (s *NoticeReceiverServiceImpl) UpdateNoticeReceiver(ctx *gin.Context,entity *models.NoticeReceiverModel) error {
    if entity.Name == "" {
        return errors.New("name is required")
    }
    return s.repo.UpdateNoticeReceiver(ctx,entity)
}

// DeleteNoticeReceiver 删除NoticeReceiver
func (s *NoticeReceiverServiceImpl) DeleteNoticeReceiver(ctx *gin.Context,id uint) error {
    return s.repo.DeleteNoticeReceiver(ctx,id)
}
// PageNoticeReceivers 分页获取NoticeReceiver列表
func (s *NoticeReceiverServiceImpl) PageNoticeReceivers(ctx *gin.Context,keywords string, pageNum, pageSize int) ([]*models.NoticeReceiverModel, int64, error) {
    return s.repo.PageNoticeReceivers(ctx,keywords, pageNum, pageSize)
}
// GetNoticeReceiverForm 表单
func (s *NoticeReceiverServiceImpl) GetNoticeReceiverForm(ctx *gin.Context,id uint) (*models.NoticeReceiverModel, error) {
    return s.repo.GetNoticeReceiverByID(ctx,id)
}
