package repositories

import (
	"errors"

	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"gorm.io/gorm"
)

// DictRepository Dict数据访问接口
type DictRepository interface {
	GetDictByID(id uint) (*models.DictModel, error) // 使用uint类型
	ListDicts() ([]*models.DictModel, error)
	CreateDict(entity *models.DictModel) error
	UpdateDict(entity *models.DictModel) error
	DeleteDict(id uint) error // 使用uint类型
	PageDicts(keywords string, pageNum, pageSize int) ([]*models.DictModel, int64, error)
	GetDictItemsByCode(dictCode string) ([]*models.DictItemModel, error)
	GetDictItemByID(id uint) (*models.DictItemModel, error)
	PageDictItems(dictCode, keywords string, pageNum, pageSize int) ([]*models.DictItemModel, int64, error)
	CreateDictItem(item *models.DictItemModel) error
	UpdateDictItem(item *models.DictItemModel) error
	DeleteDictItem(id uint) error // 使用uint类型
}

// DictRepositoryImpl Dict数据访问实现
type DictRepositoryImpl struct {
	db *gorm.DB
}

// NewDictRepository 创建Dict数据访问
func NewDictRepository(db *gorm.DB) DictRepository {
	return &DictRepositoryImpl{db: db}
}

// GetDictByID 根据ID获取Dict
func (r *DictRepositoryImpl) GetDictByID(id uint) (*models.DictModel, error) {
	var entity models.DictModel
	if err := r.db.First(&entity, id).Error; err != nil { // GORM支持直接传递uint类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// ListDicts 获取Dict列表
func (r *DictRepositoryImpl) ListDicts() ([]*models.DictModel, error) {
	var entities []*models.DictModel
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// CreateDict 创建Dict
func (r *DictRepositoryImpl) CreateDict(entity *models.DictModel) error {
	return r.db.Create(entity).Error
}

// UpdateDict 更新Dict
func (r *DictRepositoryImpl) UpdateDict(entity *models.DictModel) error {
	result := r.db.Save(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Dict not found")
	}
	return nil
}

// DeleteDict 删除Dict
func (r *DictRepositoryImpl) DeleteDict(id uint) error {
	result := r.db.Delete(&models.DictModel{}, id) // GORM支持直接传递uint类型
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Dict not found")
	}
	return nil
}

// PageDicts 分页+关键词
func (r *DictRepositoryImpl) PageDicts(keywords string, pageNum, pageSize int) ([]*models.DictModel, int64, error) {
	var list []*models.DictModel
	var total int64
	db := r.db.Model(&models.DictModel{})
	if keywords != "" {
		db = db.Where("name LIKE ?", "%"+keywords+"%")
	}
	db.Count(&total)
	err := db.Order("id desc").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// GetDictItemsByCode 根据Dict编码获取DictItem列表
func (r *DictRepositoryImpl) GetDictItemsByCode(dictCode string) ([]*models.DictItemModel, error) {
	var items []*models.DictItemModel
	if err := r.db.Where("dict_code = ?", dictCode).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetDictItemByID 根据ID获取DictItem
func (r *DictRepositoryImpl) GetDictItemByID(id uint) (*models.DictItemModel, error) {
	var item models.DictItemModel
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// ListDictItemsByDictCode 根据Dict编码获取DictItem列表
func (r *DictRepositoryImpl) ListDictItemsByDictCode(dictCode string) ([]*models.DictItemModel, error) {
	var items []*models.DictItemModel
	if err := r.db.Where("dict_code = ?", dictCode).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// PageDictItems 字典项分页实现
func (r *DictRepositoryImpl) PageDictItems(dictCode, keywords string, pageNum, pageSize int) ([]*models.DictItemModel, int64, error) {
	var list []*models.DictItemModel
	var total int64
	db := r.db.Model(&models.DictItemModel{}).Where("dict_code = ?", dictCode)
	if keywords != "" {
		db = db.Where("label LIKE ? OR value LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	db.Count(&total)
	err := db.Order("sort asc, id asc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// CreateDictItem 新增字典项
func (r *DictRepositoryImpl) CreateDictItem(item *models.DictItemModel) error {
	return r.db.Create(item).Error
}

// UpdateDictItem 更新字典项
func (r *DictRepositoryImpl) UpdateDictItem(item *models.DictItemModel) error {
	if item.ID == 0 {
		return errors.New("ID is required")
	}
	result := r.db.Save(item)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("DictItem not found")
	// }
	return nil
}

// DeleteDictItem 删除字典项
func (r *DictRepositoryImpl) DeleteDictItem(id uint) error {
	result := r.db.Delete(&models.DictItemModel{}, id) // GORM支持直接传递uint类型
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("DictItem not found")
	}
	return nil
}
