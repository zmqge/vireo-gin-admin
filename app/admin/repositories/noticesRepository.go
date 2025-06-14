package repositories

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/scopes"
	"gorm.io/gorm"
)

// NoticesRepository Notices数据访问接口
type NoticesRepository interface {
	GetNoticesByID(ctx *gin.Context, id uint) (*models.NoticesModel, error) // 使用uint类型
	ListNoticess(ctx *gin.Context) ([]*models.NoticesModel, error)
	CreateNotices(ctx *gin.Context, entity *models.NoticesModel) error
	UpdateNotices(ctx *gin.Context, entity *models.NoticesModel) error
	DeleteNotices(ctx *gin.Context, id uint) error // 使用uint类型
	PageNoticess(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error)
	RevokeNotice(ctx *gin.Context, id uint) error  // 使用uint类型
	PublishNotice(ctx *gin.Context, id uint) error // 使用uint类型

}

// NoticesRepositoryImpl Notices数据访问实现
type NoticesRepositoryImpl struct {
	db *gorm.DB
}

// NewNoticesRepository 创建Notices数据访问
func NewNoticesRepository(db *gorm.DB) NoticesRepository {
	return &NoticesRepositoryImpl{db: db}
}

// GetNoticesByID 根据ID获取Notices
func (r *NoticesRepositoryImpl) GetNoticesByID(ctx *gin.Context, id uint) (*models.NoticesModel, error) {
	var entity models.NoticesModel
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Select("notices.*, users.nickname as publisher_name").
		Joins("left join users on users.id = notices.creator_id").
		First(&entity, id).Error; err != nil { // GORM支持直接传递uint类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// ListNoticess 获取Notices列表
func (r *NoticesRepositoryImpl) ListNoticess(ctx *gin.Context) ([]*models.NoticesModel, error) {
	var entities []*models.NoticesModel
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// CreateNotices 创建Notices
func (r *NoticesRepositoryImpl) CreateNotices(ctx *gin.Context, entity *models.NoticesModel) error {
	return r.db.WithContext(context.WithValue(ctx.Request.Context(), "ginContext", ctx)).
		Create(entity).Error
}

// UpdateNotices 更新Notices
func (r *NoticesRepositoryImpl) UpdateNotices(ctx *gin.Context, entity *models.NoticesModel) error {
	result := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Save(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Notices not found")
	}
	return nil
}

// DeleteNotices 删除Notices
func (r *NoticesRepositoryImpl) DeleteNotices(ctx *gin.Context, id uint) error {
	result := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Delete(&models.NoticesModel{}, id) // GORM支持直接传递uint类型
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Notices not found")
	}
	return nil
}

// PageNoticess 分页获取Notices列表
func (r *NoticesRepositoryImpl) PageNoticess(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error) {
	var entities []*models.NoticesModel
	var total int64
	// 构建查询条件
	query := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Model(&models.NoticesModel{})
	if keywords != "" {
		query = query.Where("name LIKE ?", "%"+keywords+"%")
	}
	if publishStatus != "" {
		// 转换 publishStatus 为 uint
		status, _ := strconv.ParseUint(publishStatus, 10, 8)
		query = query.Where("status = ?", status)

	}
	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	if err := query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, total, nil
}

// RevokeNotice 撤销公告
func (r *NoticesRepositoryImpl) RevokeNotice(ctx *gin.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.NoticesModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     2,
			"revoked_at": now,
		}).Error
}

// PublishNotice 发布公告
func (r *NoticesRepositoryImpl) PublishNotice(ctx *gin.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.NoticesModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       1,
			"published_at": now,
		}).Error
}
