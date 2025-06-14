package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
)

// ConfigService Config服务接口
type ConfigService interface {
	GetConfigByID(ctx *gin.Context, id uint) (*models.ConfigModel, error) // 使用uint类型
	ListConfigs(ctx *gin.Context) ([]*models.ConfigModel, error)
	CreateConfig(ctx *gin.Context, entity *models.ConfigModel) error
	UpdateConfig(ctx *gin.Context, entity *models.ConfigModel) error
	DeleteConfig(ctx *gin.Context, id uint) error // 使用uint类型
	PageConfigs(ctx *gin.Context, keywords string, pageNum, pageSize int) ([]*models.ConfigModel, int64, error)
	GetConfigForm(ctx *gin.Context, id uint) (*models.ConfigModel, error) // 使用uint类型
}

// ConfigServiceImpl Config服务实现
type ConfigServiceImpl struct {
	repo repositories.ConfigRepository
}

// NewConfigService 创建Config服务
func NewConfigService(repo repositories.ConfigRepository) ConfigService {
	return &ConfigServiceImpl{repo: repo}
}

// GetConfigByID 根据ID获取Config
func (s *ConfigServiceImpl) GetConfigByID(ctx *gin.Context, id uint) (*models.ConfigModel, error) {
	entity, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("Config not found")
	}
	return entity, nil
}

// ListConfigs 获取Config列表
func (s *ConfigServiceImpl) ListConfigs(ctx *gin.Context) ([]*models.ConfigModel, error) {
	return s.repo.ListConfigs(ctx)
}

// CreateConfig 创建Config
func (s *ConfigServiceImpl) CreateConfig(ctx *gin.Context, entity *models.ConfigModel) error {
	if entity.ConfigName == "" {
		return errors.New("name is required")
	}
	return s.repo.CreateConfig(ctx, entity)
}

// UpdateConfig 更新Config
func (s *ConfigServiceImpl) UpdateConfig(ctx *gin.Context, entity *models.ConfigModel) error {
	if entity.ConfigName == "" {
		return errors.New("name is required")
	}
	return s.repo.UpdateConfig(ctx, entity)
}

// DeleteConfig 删除Config
func (s *ConfigServiceImpl) DeleteConfig(ctx *gin.Context, id uint) error {
	return s.repo.DeleteConfig(ctx, id)
}

// PageConfigs 分页获取Config列表
func (s *ConfigServiceImpl) PageConfigs(ctx *gin.Context, keywords string, pageNum, pageSize int) ([]*models.ConfigModel, int64, error) {
	return s.repo.PageConfigs(ctx, keywords, pageNum, pageSize)
}

// GetConfigForm 表单
func (s *ConfigServiceImpl) GetConfigForm(ctx *gin.Context, id uint) (*models.ConfigModel, error) {
	return s.repo.GetConfigByID(ctx, id)
}
