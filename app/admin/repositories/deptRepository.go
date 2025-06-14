package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/scopes"
	"gorm.io/gorm"
)

type DeptRepository struct {
	db *gorm.DB
}

func NewDeptRepository(db *gorm.DB) *DeptRepository {
	return &DeptRepository{db: db}
}

func (r *DeptRepository) GetDeptByID(id uint) (*models.Dept, error) {
	var dept models.Dept
	if err := r.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DeptRepository) ListDepts(ctx *gin.Context) ([]models.Dept, error) {
	var depts []models.Dept
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Order("sort ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *DeptRepository) CreateDept(dept *models.Dept) error {
	return r.db.Create(dept).Error
}

func (r *DeptRepository) UpdateDept(dept *models.Dept) error {
	updateMap := map[string]interface{}{
		"parent_id": dept.ParentID,
		"name":      dept.Name,
		"code":      dept.Code,
		"sort":      dept.Sort,
		"status":    dept.Status,
	}
	return r.db.Model(&models.Dept{}).Where("id = ?", dept.ID).Updates(updateMap).Error
}

func (r *DeptRepository) DeleteDept(id uint) error {
	return r.db.Delete(&models.Dept{}, id).Error
}
