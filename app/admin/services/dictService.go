package services

import (
	"errors"

	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
)

// DictService Dict服务接口
type DictService interface {
	GetDictByID(id uint) (*models.DictModel, error) // 使用uint类型
	ListDicts() ([]*models.DictModel, error)
	CreateDict(entity *models.DictModel) error
	UpdateDict(entity *models.DictModel) error
	DeleteDict(id uint) error // 使用uint类型
	PageDicts(keywords string, pageNum, pageSize int) ([]*models.DictModel, int64, error)
	GetDictItemsByCode(dictCode string) ([]*models.DictItemModel, error)
	PageDictItems(dictCode, keywords string, pageNum, pageSize int) ([]*models.DictItemModel, int64, error)
	CreateDictItem(item *models.DictItemModel) error
	UpdateDictItem(item *models.DictItemModel) error
	DeleteDictItem(id uint) error                                                // 使用uint类型
	GetDictItemForm(dictCode string, itemId uint) (*models.DictItemModel, error) // 新增
}

// DictServiceImpl Dict服务实现
type DictServiceImpl struct {
	repo repositories.DictRepository
}

// NewDictService 创建Dict服务
func NewDictService(repo repositories.DictRepository) DictService {
	return &DictServiceImpl{repo: repo}
}

// GetDictByID 根据ID获取Dict
func (s *DictServiceImpl) GetDictByID(id uint) (*models.DictModel, error) {
	entity, err := s.repo.GetDictByID(id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("Dict not found")
	}
	return entity, nil
}

// ListDicts 获取Dict列表
func (s *DictServiceImpl) ListDicts() ([]*models.DictModel, error) {
	return s.repo.ListDicts()
}

// CreateDict 创建Dict
func (s *DictServiceImpl) CreateDict(entity *models.DictModel) error {
	if entity.Name == "" {
		return errors.New("name is required")
	}
	return s.repo.CreateDict(entity)
}

// UpdateDict 更新Dict
func (s *DictServiceImpl) UpdateDict(entity *models.DictModel) error {
	if entity.Name == "" {
		return errors.New("name is required")
	}
	return s.repo.UpdateDict(entity)
}

// DeleteDict 删除Dict
func (s *DictServiceImpl) DeleteDict(id uint) error {
	return s.repo.DeleteDict(id)
}

// PageDicts 分页查询
func (s *DictServiceImpl) PageDicts(keywords string, pageNum, pageSize int) ([]*models.DictModel, int64, error) {
	return s.repo.PageDicts(keywords, pageNum, pageSize)
}

// GetDictItemsByCode 根据Dict编码获取DictItem列表
func (s *DictServiceImpl) GetDictItemsByCode(dictCode string) ([]*models.DictItemModel, error) {
	return s.repo.GetDictItemsByCode(dictCode)
}

// PageDictItems 字典项分页
func (s *DictServiceImpl) PageDictItems(dictCode, keywords string, pageNum, pageSize int) ([]*models.DictItemModel, int64, error) {
	return s.repo.PageDictItems(dictCode, keywords, pageNum, pageSize)
}

// CreateDictItem 新增字典项
func (s *DictServiceImpl) CreateDictItem(item *models.DictItemModel) error {
	return s.repo.CreateDictItem(item)
}

// UpdateDictItem 更新字典项
func (s *DictServiceImpl) UpdateDictItem(item *models.DictItemModel) error {
	if item.ID == 0 {
		return errors.New("ID is required")
	}
	if item.DictCode == "" {
		return errors.New("dictCode is required")
	}
	if item.Value == "" {
		return errors.New("value is required")
	}
	if item.Label == "" {
		return errors.New("label is required")
	}
	return s.repo.UpdateDictItem(item)
}

// 删除DictItem
func (s *DictServiceImpl) DeleteDictItem(id uint) error {
	if id == 0 {
		return errors.New("ID is required")
	}
	return s.repo.DeleteDictItem(id)
}

// GetDictItemForm 获取字典项表单数据
func (s *DictServiceImpl) GetDictItemForm(dictCode string, itemId uint) (*models.DictItemModel, error) {
	item, err := s.repo.GetDictItemByID(itemId)
	if err != nil {
		return nil, err
	}
	if item == nil || item.DictCode != dictCode {
		return nil, nil
	}
	return item, nil
}
