package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/scopes"
	"gorm.io/gorm"
)

// NoticeReceiverRepository NoticeReceiver数据访问接口
type NoticeReceiverRepository interface {
	GetNoticeReceiverByID(ctx *gin.Context, id uint) (*models.NoticeReceiverModel, error) // 使用uint类型
	ListNoticeReceivers(ctx *gin.Context) ([]*models.NoticeReceiverModel, error)
	CreateNoticeReceiver(ctx *gin.Context, entity *models.NoticeReceiverModel) error
	UpdateNoticeReceiver(ctx *gin.Context, entity *models.NoticeReceiverModel) error
	DeleteNoticeReceiver(ctx *gin.Context, id uint) error // 使用uint类型
	PageNoticeReceivers(ctx *gin.Context, keywords string, pageNum, pageSize int) ([]*models.NoticeReceiverModel, int64, error)

	BatchCreate(ctx *gin.Context, receivers []models.NoticeReceiver) error
	BatchCreateTx(ctx *gin.Context, tx *gorm.DB, receivers []models.NoticeReceiver) error // 新增事务版本
}

// NoticeReceiverRepositoryImpl NoticeReceiver数据访问实现
type NoticeReceiverRepositoryImpl struct {
	db *gorm.DB
}

// NewNoticeReceiverRepository 创建NoticeReceiver数据访问
func NewNoticeReceiverRepository(db *gorm.DB) NoticeReceiverRepository {
	return &NoticeReceiverRepositoryImpl{db: db}
}

// GetNoticeReceiverByID 根据ID获取NoticeReceiver
func (r *NoticeReceiverRepositoryImpl) GetNoticeReceiverByID(ctx *gin.Context, id uint) (*models.NoticeReceiverModel, error) {
	var entity models.NoticeReceiverModel
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		First(&entity, id).Error; err != nil { // GORM支持直接传递uint类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// ListNoticeReceivers 获取NoticeReceiver列表
func (r *NoticeReceiverRepositoryImpl) ListNoticeReceivers(ctx *gin.Context) ([]*models.NoticeReceiverModel, error) {
	var entities []*models.NoticeReceiverModel
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// CreateNoticeReceiver 创建NoticeReceiver
func (r *NoticeReceiverRepositoryImpl) CreateNoticeReceiver(ctx *gin.Context, entity *models.NoticeReceiverModel) error {
	return r.db.WithContext(context.WithValue(ctx.Request.Context(), "ginContext", ctx)).
		Create(entity).Error
}

// UpdateNoticeReceiver 更新NoticeReceiver
func (r *NoticeReceiverRepositoryImpl) UpdateNoticeReceiver(ctx *gin.Context, entity *models.NoticeReceiverModel) error {
	result := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Save(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("NoticeReceiver not found")
	}
	return nil
}

// DeleteNoticeReceiver 删除NoticeReceiver
func (r *NoticeReceiverRepositoryImpl) DeleteNoticeReceiver(ctx *gin.Context, id uint) error {
	result := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Delete(&models.NoticeReceiverModel{}, id) // GORM支持直接传递uint类型
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("NoticeReceiver not found")
	}
	return nil
}

// PageNoticeReceivers 分页获取NoticeReceiver列表
func (r *NoticeReceiverRepositoryImpl) PageNoticeReceivers(ctx *gin.Context, keywords string, pageNum, pageSize int) ([]*models.NoticeReceiverModel, int64, error) {
	var entities []*models.NoticeReceiverModel
	var total int64
	// 构建查询条件
	query := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Model(&models.NoticeReceiverModel{})
	if keywords != "" {
		query = query.Where("name LIKE ?", "%"+keywords+"%")
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

func (r *NoticeReceiverRepositoryImpl) BatchCreate(ctx *gin.Context, receivers []models.NoticeReceiver) error {
	if len(receivers) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(receivers, 500).Error
}

func (r *NoticeReceiverRepositoryImpl) BatchCreateTx(ctx *gin.Context, tx *gorm.DB, receivers []models.NoticeReceiver) error {
	log.Printf("准备批量插入 %d 条接收者记录\n", len(receivers))

	if len(receivers) == 0 {
		return nil
	}

	return tx.CreateInBatches(receivers, 500).Error
}
