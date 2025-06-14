package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
	"gorm.io/gorm"
)

const (
	maxTreeDepth    = 10
	deptCacheKey    = "dept:all"
	cacheExpiration = 5 * time.Minute
)

type DeptService struct {
	repo  *repositories.DeptRepository
	cache *cache.Cache // 实际使用时需要初始化
}

func NewDeptService(db *gorm.DB) *DeptService {
	return &DeptService{
		repo:  repositories.NewDeptRepository(db),
		cache: cache.New(cacheExpiration, 10*time.Minute),
	}
}

// CreateDept 创建部门
func (s *DeptService) CreateDept(dept *models.Dept) error {
	if err := s.repo.CreateDept(dept); err != nil {
		return err
	}
	s.cache.Delete(deptCacheKey)
	return nil
}

// UpdateDept 更新部门信息
func (s *DeptService) UpdateDept(ctx *gin.Context, dept *models.Dept) error {
	if dept.ParentID == dept.ID {
		return errors.New("上级部门不能是本部门")
	}
	// 检查部门是否存在
	if _, err := s.repo.GetDeptByID(dept.ID); err != nil {
		return errors.New("部门不存在")
	}
	// 检查父部门是否存在（如果ParentID不为0）
	if dept.ParentID != 0 {
		if _, err := s.repo.GetDeptByID(dept.ParentID); err != nil {
			return errors.New("指定的上级部门不存在")
		}
	}
	if err := s.checkCircularReference(ctx, dept.ID, dept.ParentID); err != nil {
		return err
	}
	if err := s.repo.UpdateDept(dept); err != nil {
		return err
	}
	s.cache.Delete(deptCacheKey)
	return nil
}

// checkCircularReference 检查部门修改是否会造成循环引用
func (s *DeptService) checkCircularReference(ctx *gin.Context, deptID, newParentID uint) error {
	if newParentID == 0 {
		return nil
	}
	depts, err := s.getAllDepts(ctx)
	if err != nil {
		return err
	}
	childToParent := make(map[uint]uint)
	for _, d := range depts {
		if d.ID == deptID {
			continue
		}
		childToParent[d.ID] = d.ParentID
	}
	currentParent := newParentID
	for {
		if currentParent == deptID {
			return errors.New("修改会导致循环引用：指定的上级部门已经是本部门的子部门")
		}
		if currentParent == 0 {
			break
		}
		nextParent, exists := childToParent[currentParent]
		if !exists {
			break
		}
		currentParent = nextParent
	}
	return nil
}

// GetDeptDetails 获取部门详情
func (s *DeptService) GetDeptDetails(id uint) (*models.Dept, error) {
	return s.repo.GetDeptByID(id)
}

// 基础数据获取方法
func (s *DeptService) getAllDepts(ctx *gin.Context) ([]models.Dept, error) {
	if cached, found := s.cache.Get(deptCacheKey); found {
		return cached.([]models.Dept), nil
	}
	depts, err := s.repo.ListDepts(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set(deptCacheKey, depts, cacheExpiration)
	return depts, nil
}

// GetDeptOptions 获取部门选项（用于下拉选择等场景）
func (s *DeptService) GetDeptOptions(ctx *gin.Context) ([]models.OptionLong, error) {
	depts, err := s.getAllDepts(ctx)
	if err != nil {
		return nil, err
	}
	return BuildDeptOptions(depts, 0), nil
}

// GetDepts 获取部门列表（带过滤条件）
func (s *DeptService) GetDepts(ctx *gin.Context, keywords string, status string) ([]models.DeptV0, error) {
	depts, err := s.repo.ListDepts(ctx)
	if err != nil {
		return nil, err
	}
	var filtered []models.Dept
	for _, d := range depts {
		if keywords != "" && !contains(d.Name, keywords) {
			continue
		}
		if status != "" && fmt.Sprintf("%d", d.Status) != status {
			continue
		}
		filtered = append(filtered, d)
	}
	return BuildDeptTree(filtered)
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr)))))
}

// DeleteDept 删除部门
func (s *DeptService) DeleteDept(ctx *gin.Context, id uint) error {
	if _, err := s.repo.GetDeptByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("部门不存在")
		}
		return fmt.Errorf("查询部门失败: %v", err)
	}
	depts, err := s.repo.ListDepts(ctx)
	if err != nil {
		return fmt.Errorf("查询子部门失败: %v", err)
	}
	for _, d := range depts {
		if d.ParentID == id {
			return errors.New("该部门下有子部门，请先删除或转移子部门")
		}
	}
	// 检查是否有用户关联（此处建议在repo层实现更优）
	// ...如有UserRepository可调用...
	if err := s.repo.DeleteDept(id); err != nil {
		return fmt.Errorf("删除部门失败: %v", err)
	}
	s.cache.Delete(deptCacheKey)
	return nil
}

// BuildDeptOptions 构建部门选项树
func BuildDeptOptions(depts []models.Dept, parentID uint) []models.OptionLong {
	childMap := buildParentChildMap(depts)
	return buildOptionsFromMap(childMap, parentID)
}

func buildOptionsFromMap(childMap map[uint][]models.Dept, parentID uint) []models.OptionLong {
	var options []models.OptionLong

	for _, dept := range childMap[parentID] {
		option := models.OptionLong{
			Value: dept.ID,
			Label: dept.Name,
		}

		if children := buildOptionsFromMap(childMap, dept.ID); len(children) > 0 {
			option.Children = children
		}

		options = append(options, option)
	}

	return options
}

// BuildDeptTree 构建部门树
func BuildDeptTree(depts []models.Dept) ([]models.DeptV0, error) {
	childMap := buildParentChildMap(depts)
	rootNodes := findRootNodes(depts)

	var result []models.DeptV0
	for _, rootID := range rootNodes {
		tree, err := buildTreeFromMap(childMap, rootID, 0)
		if err != nil {
			return nil, err
		}
		result = append(result, tree...)
	}

	return result, nil
}

func buildTreeFromMap(childMap map[uint][]models.Dept, parentID uint, depth int) ([]models.DeptV0, error) {
	if depth > maxTreeDepth {
		return nil, errors.New("部门层级过深，可能存在循环引用")
	}

	var result []models.DeptV0
	for _, dept := range childMap[parentID] {
		deptV0 := models.DeptV0{
			ID:       dept.ID,
			ParentID: dept.ParentID,
			Name:     dept.Name,
			Code:     dept.Code,
			Sort:     dept.Sort,
			Status:   dept.Status,
		}

		children, err := buildTreeFromMap(childMap, dept.ID, depth+1)
		if err != nil {
			return nil, err
		}

		if len(children) > 0 {
			deptV0.Children = &children
		}

		result = append(result, deptV0)
	}
	return result, nil
}

// 辅助函数
func buildParentChildMap(depts []models.Dept) map[uint][]models.Dept {
	m := make(map[uint][]models.Dept)
	for _, d := range depts {
		m[d.ParentID] = append(m[d.ParentID], d)
	}
	return m
}

func findRootNodes(depts []models.Dept) []uint {
	allDeptIDs := make(map[uint]bool)
	for _, d := range depts {
		allDeptIDs[d.ID] = true
	}

	rootSet := make(map[uint]bool)
	for _, d := range depts {
		if !allDeptIDs[d.ParentID] || d.ParentID == 0 {
			rootSet[d.ParentID] = true
		}
	}

	var roots []uint
	for rootID := range rootSet {
		roots = append(roots, rootID)
	}
	return roots
}
