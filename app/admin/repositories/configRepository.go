package repositories

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/scopes"
	"gorm.io/gorm"
)

// ConfigRepository Config数据访问接口
type ConfigRepository interface {
	GetConfigByID(ctx *gin.Context, id uint) (*models.ConfigModel, error) // 使用uint类型
	ListConfigs(ctx *gin.Context) ([]*models.ConfigModel, error)
	CreateConfig(ctx *gin.Context, entity *models.ConfigModel) error
	UpdateConfig(ctx *gin.Context, entity *models.ConfigModel) error
	DeleteConfig(ctx *gin.Context, id uint) error // 使用uint类型
	PageConfigs(ctx *gin.Context, keywords string, pageNum, pageSize int) ([]*models.ConfigModel, int64, error)
}

// ConfigRepositoryImpl Config数据访问实现
type ConfigRepositoryImpl struct {
	db *gorm.DB
}

// NewConfigRepository 创建Config数据访问
func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &ConfigRepositoryImpl{db: db}
}

// GetConfigByID 根据ID获取Config
func (r *ConfigRepositoryImpl) GetConfigByID(ctx *gin.Context, id uint) (*models.ConfigModel, error) {
	var entity models.ConfigModel
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		First(&entity, id).Error; err != nil { // GORM支持直接传递uint类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// ListConfigs 获取Config列表
func (r *ConfigRepositoryImpl) ListConfigs(ctx *gin.Context) ([]*models.ConfigModel, error) {
	var entities []*models.ConfigModel
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// CreateConfig 创建Config
func (r *ConfigRepositoryImpl) CreateConfig(ctx *gin.Context, entity *models.ConfigModel) error {
	return r.db.WithContext(context.WithValue(ctx.Request.Context(), "ginContext", ctx)).
		Create(entity).Error
}

// UpdateConfig 更新Config
func (r *ConfigRepositoryImpl) UpdateConfig(ctx *gin.Context, entity *models.ConfigModel) error {
	result := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Save(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Config not found")
	}
	return nil
}

// DeleteConfig 删除Config
func (r *ConfigRepositoryImpl) DeleteConfig(ctx *gin.Context, id uint) error {
	result := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Delete(&models.ConfigModel{}, id) // GORM支持直接传递uint类型
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Config not found")
	}
	return nil
}

// PageConfigs 分页获取Config列表
func (r *ConfigRepositoryImpl) PageConfigs(ctx *gin.Context, keywords string, pageNum, pageSize int) ([]*models.ConfigModel, int64, error) {
	var entities []*models.ConfigModel
	var total int64
	// 构建查询条件
	query := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Model(&models.ConfigModel{})
	if keywords != "" {
		query = query.Where("config_name LIKE ?", "%"+keywords+"%")
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
