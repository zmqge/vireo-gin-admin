package repositories

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/scopes"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	// 数据库操作：查询用户
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	// 数据库操作：根据用户名查询用户
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Create(user *models.User) error {
	// 数据库操作：创建用户记录
	return r.db.Create(user).Error
}

// 删除用户（软删除，用户名追加 _deleted_时间戳）
func (r *UserRepository) Delete(id uint) error {
	// 先查出用户
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}
	// 修改用户名，防止唯一索引冲突，追加 _deleted_加时间戳
	timestamp := time.Now().Format("20060102150405")
	newUsername := user.Username + "_deleted_" + timestamp
	if err := r.db.Model(&user).Update("username", newUsername).Error; err != nil {
		return err
	}
	// 执行软删除
	return r.db.Delete(&user).Error
}

func (r *UserRepository) GetUserPage(ctx *gin.Context, params models.UserQueryParams) (*models.UserPageResult, error) {
	// 数据库操作：获取用户分页数据
	query := r.db.Table("users as u").
		Scopes(scopes.DataPermissionScope(ctx)).WithContext(ctx).
		Model(&models.User{})

	// 数据库操作：应用过滤条件
	applyFilters(query, params)

	// 数据库操作：获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// 数据库操作：获取分页数据
	var users []models.User
	err := query.
		Offset((params.PageNum - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&users).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return &models.UserPageResult{
		Users: users,
		Total: total,
	}, nil
}

// 辅助函数：应用过滤条件
func applyFilters(query *gorm.DB, params models.UserQueryParams) {
	// 关键字搜索
	if params.Keywords != "" {
		searchPattern := "%" + params.Keywords + "%"
		query = query.Where("u.username LIKE ? OR u.nickname LIKE ? OR u.mobile LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 状态过滤
	if params.Status != "" {
		query = query.Where("u.status = ?", params.Status)
	}

	// 部门过滤
	if params.DeptID > 0 {
		query = query.Where("u.dept_id = ?", params.DeptID)
	}

	// 创建时间范围过滤
	// 创建时间范围过滤
	if len(params.CreateTime) == 2 && params.CreateTime[0] != "" && params.CreateTime[1] != "" {
		startTime, err := time.Parse("2006-01-02", params.CreateTime[0])
		if err != nil {
			// log or handle error as needed, but just return
			return
		}
		startTime = startTime.UTC()

		endTime, err := time.Parse("2006-01-02", params.CreateTime[1])
		if err != nil {
			// log or handle error as needed, but just return
			return
		}
		endTime = endTime.UTC().Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		query = query.Where("u.created_at BETWEEN ? AND ?", startTime, endTime)
	}

	// 角色ID过滤
	if len(params.RoleIDs) > 0 && params.RoleIDs[0] != "" {
		query = query.
			Joins("JOIN user_roles ur ON ur.user_id = u.id").
			Where("ur.role_id IN ?", params.RoleIDs)
	}
}

// 辅助函数：限制值在[min,max]范围内，否则返回defaultValue
func clamp(value, min, max, defaultValue int) int {
	if value <= 0 {
		return defaultValue
	}
	if value > max {
		return max
	}
	return value
}

// 辅助函数：返回两个数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// UpdateUserFull 全量更新用户信息及角色
func (r *UserRepository) UpdateUserFull(id uint, nickname, mobile string, gender string, avatar, email string, status, deptId int, roleIds []int64, openId string) error {
	updates := map[string]interface{}{
		"nickname": nickname,
		"mobile":   mobile,
		"gender":   gender,
		"avatar":   avatar,
		"email":    email,
		"status":   status,
		"dept_id":  deptId,
	}
	// openId 字段如有可加: "open_id": openId
	// 数据库操作：更新用户信息
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}

	// 数据库操作：更新用户角色关联
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}
	if err := r.db.Model(&user).Association("RoleList").Clear(); err != nil {
		return err
	}
	if len(roleIds) > 0 {
		var roles []models.Role
		for _, rid := range roleIds {
			roles = append(roles, models.Role{ID: uint(rid)})
		}
		// 数据库操作：替换用户角色
		if err := r.db.Model(&user).Association("RoleList").Replace(roles); err != nil {
			return err
		}
	}
	return nil
}

// CreateUserFull 创建用户及角色
func (r *UserRepository) CreateUserFull(username, nickname, mobile string, gender string, avatar, email string, status int, deptId uint, roleIds []int64, openId, password string) error {
	// 生成 salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	saltStr := base64.StdEncoding.EncodeToString(salt)
	// 密码加密：password+salt
	hash, err := bcrypt.GenerateFromPassword([]byte(password+saltStr), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		Username: username,
		Nickname: nickname,
		Mobile:   mobile,
		Gender:   gender,
		Avatar:   avatar,
		Email:    email,
		Status:   status,
		DeptID:   deptId,
		Password: string(hash),
		Salt:     saltStr,
		// OpenId 字段如有可加: OpenId: openId
	}
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	if len(roleIds) > 0 {
		var roles []models.Role
		for _, rid := range roleIds {
			roles = append(roles, models.Role{ID: uint(rid)})
		}
		if err := r.db.Model(&user).Association("RoleList").Replace(roles); err != nil {
			return err
		}
	}
	return nil
}

// UpdatePassword 重置用户密码（含salt）
func (r *UserRepository) UpdatePassword(id uint, password, salt string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password": password,
		"salt":     salt,
	}).Error
}

// UpdateUserProfile 更新个人中心用户信息（仅允许部分字段，map 更新，支持零值）
func (r *UserRepository) UpdateUserProfile(id uint, updateMap map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updateMap).Error
}

// 可根据需要继续补充 Update 等方法

func (r *UserRepository) ListUserOptions(ctx *gin.Context) ([]models.UserOption, error) {
	var users []models.User
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).WithContext(ctx).
		Select("id , username ").
		Find(&users).Error; err != nil {
		return []models.UserOption{}, err
	}

	options := make([]models.UserOption, len(users))
	for i, user := range users {
		// 修复 UserOption 未定义的问题，使用正确的类型 models.UserOption
		options[i] = models.UserOption{
			Value: user.ID,
			Label: user.Username,
		}
	}
	// 修复 err 未定义的问题，使用查询时的错误检查
	// 由于在之前的 Find 方法调用中已经有错误检查，这里复用该错误变量
	if err := r.db.Scopes(scopes.DataPermissionScope(ctx)).
		Table("users").
		Find(&users).Error; err != nil {
		return []models.UserOption{}, err
	}
	return options, nil
}
