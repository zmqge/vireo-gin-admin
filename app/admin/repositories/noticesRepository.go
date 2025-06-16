package repositories

import (
	"context"
	"errors"
	"fmt"
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
	PageNotices(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error)
	RevokeNotice(ctx *gin.Context, id uint) error  // 使用uint类型
	PublishNotice(ctx *gin.Context, id uint) error // 使用uint类型
	MarkNoticeAsRead(ctx *gin.Context, userID uint, noticeID uint) error

	GetNoticeWithReceivers(ctx *gin.Context, id uint) (*models.NoticesModel, error)
	UpdateNoticeStatusTx(tx *gorm.DB, id uint, status uint) error

	// 新增事务支持
	CreateNoticesTx(tx *gorm.DB, entity *models.NoticesModel) error
	BeginTx(ctx *gin.Context) *gorm.DB
	PublishNoticesTx(ctx *gin.Context, tx *gorm.DB, id uint) error
	RevokeNoticesTx(tx *gorm.DB, id uint) error
	MyPageNotices(ctx *gin.Context, userID uint, keywords string, isRead uint, pageNum, pageSize int) ([]*models.NoticesModel, int64, error)
	GetMyNoticesByID(ctx *gin.Context, userID uint, noticeID uint) (*models.NoticesModel, error)
	MarkAllAsRead(ctx *gin.Context, userID uint) error
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

// GetMyNoticesByID 根据用户ID和通知ID获取用户的通知
func (r *NoticesRepositoryImpl) GetMyNoticesByID(ctx *gin.Context, userID uint, noticeID uint) (*models.NoticesModel, error) {
	var notice models.NoticesModel

	// 构建查询：通过notice_receiver表关联，确保用户有权限访问该通知
	err := r.db.WithContext(ctx).
		Select("notices.*, users.nickname as publisher_name, notice_receiver.is_read").
		Joins("JOIN notice_receiver ON notices.id = notice_receiver.notice_id").
		Joins("LEFT JOIN users ON users.id = notices.creator_id").
		Where("notice_receiver.user_id = ? AND notice_receiver.notice_id = ?", userID, noticeID).
		First(&notice).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 无记录不算错误
		}
		return nil, fmt.Errorf("查询用户通知失败: %w", err)
	}

	return &notice, nil
}

// MarkNoticeAsRead 标记通知为已读
func (r *NoticesRepositoryImpl) MarkNoticeAsRead(ctx *gin.Context, userID uint, noticeID uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.NoticeReceiver{}).
		Where("user_id = ? AND notice_id = ?", userID, noticeID).
		Updates(map[string]interface{}{
			"is_read":   1,
			"read_time": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("更新已读状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到对应的通知接收记录")
	}

	return nil
}

// MarkAllAsRead 标记用户所有通知为已读
func (r *NoticesRepositoryImpl) MarkAllAsRead(ctx *gin.Context, userID uint) error {
	// 使用批量更新语法（MySQL示例）
	err := r.db.WithContext(ctx).Exec(`
        UPDATE notice_receiver 
        SET is_read = 1, read_time = NOW() 
        WHERE user_id = ? AND is_read = 0
    `, userID).Error

	if err != nil {
		return fmt.Errorf("更新所有通知为已读状态失败: %w", err)
	}
	return nil
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
func (r *NoticesRepositoryImpl) PageNotices(ctx *gin.Context, keywords string, publishStatus string, pageNum, pageSize int) ([]*models.NoticesModel, int64, error) {
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

// CreateNoticesTx 创建事务
func (r *NoticesRepositoryImpl) CreateNoticesTx(tx *gorm.DB, entity *models.NoticesModel) error {
	return tx.Create(entity).Error
}

// BeginTx 开启事务
func (r *NoticesRepositoryImpl) BeginTx(ctx *gin.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

// PublishNoticesTx 发布事务
func (r *NoticesRepositoryImpl) PublishNoticesTx(ctx *gin.Context, tx *gorm.DB, id uint) error {
	now := time.Now()
	return tx.Model(&models.NoticesModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       1,
			"published_at": now,
		}).Error
}

// RevokeNoticesTx 撤销事务
func (r *NoticesRepositoryImpl) RevokeNoticesTx(tx *gorm.DB, id uint) error {
	now := time.Now()
	return tx.Model(&models.NoticesModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     2,
			"revoked_at": now,
		}).Error
}

// MyPageNotices 获取我的公告列表
// MyPageNotices 获取我的公告列表（带已读/未读状态）
func (r *NoticesRepositoryImpl) MyPageNotices(ctx *gin.Context, userID uint, keywords string, isRead uint, pageNum, pageSize int) ([]*models.NoticesModel, int64, error) {
	var notices []*models.NoticesModel
	var total int64

	// 构建基础查询
	query := r.db.WithContext(ctx).
		Model(&models.NoticesModel{}).
		Select("notices.*, notice_receiver.is_read").
		Joins("JOIN notice_receiver ON notices.id = notice_receiver.notice_id").
		Where("notice_receiver.user_id = ?", userID)

	// 添加关键词筛选
	if keywords != "" {
		query = query.Where("notices.title LIKE ?", "%"+keywords+"%")
	}

	// 添加已读/未读筛选
	if isRead == 0 || isRead == 1 {
		query = query.Where("notice_receiver.is_read = ?", isRead)
	}
	// 必须是已发布的消息，其他的状态不展示给用户
	query = query.Where("notices.status = 1")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取总数失败: %w", err)
	}

	// 执行分页查询
	err := query.
		Order("notices.id DESC").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&notices).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询通知列表失败: %w", err)
	}

	return notices, total, nil
}

// 获取通知及其接收者关系
func (r *NoticesRepositoryImpl) GetNoticeWithReceivers(ctx *gin.Context, id uint) (*models.NoticesModel, error) {
	var notice models.NoticesModel
	err := r.db.WithContext(ctx).
		Preload("Receivers").
		First(&notice, id).Error
	if err != nil {
		return nil, err
	}
	return &notice, nil
}

// 更新通知状态（事务版本）
func (r *NoticesRepositoryImpl) UpdateNoticeStatusTx(tx *gorm.DB, id uint, status uint) error {
	return tx.Model(&models.NoticesModel{}).
		Where("id = ?", id).
		Update("status", status).Error
}
