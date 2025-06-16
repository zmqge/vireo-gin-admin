package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/repositories"
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"github.com/zmqge/vireo-gin-admin/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetList(page, pageSize int) ([]models.User, int64, error)
	CreateUser(username, password string) error
	UpdateUser(id string, username string, status int) error
	UpdateUserFull(id string, nickname, mobile string, gender string, avatar, email string, status, deptId int, roleIds []int64, openId string) error
	CreateUserFull(username, nickname, mobile string, gender string, avatar, email string, status int, deptId uint, roleIds []int64, openId, password string) error
	ChangePassword(userID, oldPassword, newPassword string) error
	GetDeptName(deptID uint) (string, error)
	GetRoleNames(userID string) (string, error)
	UpdateUserProfile(userID string, nickname, avatar string, gender string, mobile, email string) error
	ListUserOptions(ctx *gin.Context) ([]models.UserOption, error)
	Delete(userID string) error
	GetUser(userID string) (*models.User, error)
	GetUserPage(ctx *gin.Context, params models.UserQueryParams) (*models.UserPageResult, error)
	GetUserRoles(userID string) ([]models.Role, error)
	ResetPassword(userID string, password string) error
	VerifyUser(username, password string) (*models.User, error)
	UpdateLastLogin(userID uint, ClientIP string, loginTime time.Time, userAgent string) error
}
type UserServiceImpl struct {
	db   *gorm.DB
	repo repositories.UserRepository
}

func NewUserService(db *gorm.DB) *UserServiceImpl {
	return &UserServiceImpl{db: db, repo: repositories.NewUserRepository(db)}
}

func (s *UserServiceImpl) GetList(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if s.db == nil {
		return nil, 0, fmt.Errorf("database connection is nil")
	}

	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// 生成随机 salt
func RandSalt() string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	return base64.StdEncoding.EncodeToString(salt)
}

func (s *UserServiceImpl) CreateUser(username, password string) error {
	salt := RandSalt()
	hashedPassword, err := models.HashPasswordWithSalt(password, salt)
	if err != nil {
		return err
	}
	user := models.User{
		Username: username,
		Password: hashedPassword,
		Salt:     salt,
	}
	return s.repo.Create(&user)
}

func (s *UserServiceImpl) VerifyUser(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	// 校验状态
	if user.Status != 1 {
		return nil, fmt.Errorf("用户已被禁用")
	}
	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.Salt)); err != nil {
		return nil, err
	}

	return &user, nil
}

// 全局函数（兼容旧代码）
func GetUserList() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(username, password string) error {
	return NewUserService(database.DB).CreateUser(username, password)
}

func (s *UserServiceImpl) Delete(id string) error {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(uint(uid))
}

func (s *UserServiceImpl) UpdateUser(id string, username string, status int) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"username": username,
			"status":   status,
		}).Error
}

// UpdateUserFull 全量更新用户信息及角色
func (s *UserServiceImpl) UpdateUserFull(id string, nickname, mobile string, gender string, avatar, email string, status, deptId int, roleIds []int64, openId string) error {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserFull(uint(uid), nickname, mobile, gender, avatar, email, status, deptId, roleIds, openId)
}

// GetUser 根据ID获取用户
func (s *UserServiceImpl) GetUser(id string) (*models.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(uint(uid))
}

// GetUserPage 分页查询用户
func (s *UserServiceImpl) GetUserPage(ctx *gin.Context, params models.UserQueryParams) (*models.UserPageResult, error) {
	// 这里建议将分页查询下沉到 repo 层，service 只做业务逻辑
	return s.repo.GetUserPage(ctx, params)
}

// GetUserRoles 获取用户角色列表
func (s *UserServiceImpl) GetUserRoles(userID string) ([]models.Role, error) {
	var user models.User
	if err := s.db.Preload("RoleList").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.RoleList, nil
}

// CreateUserFull 创建用户及角色
func (s *UserServiceImpl) CreateUserFull(username, nickname, mobile string, gender string, avatar, email string, status int, deptId uint, roleIds []int64, openId, password string) error {
	return s.repo.CreateUserFull(username, nickname, mobile, gender, avatar, email, status, deptId, roleIds, openId, password)
}

// ResetPassword 重置用户密码
func (s *UserServiceImpl) ResetPassword(userID string, password string) error {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	salt := RandSalt()
	hash, err := models.HashPasswordWithSalt(password, salt)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(uint(uid), hash, salt)
}

// ChangePassword 修改当前用户密码，校验原密码
func (s *UserServiceImpl) ChangePassword(userID, oldPassword, newPassword string) error {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	user, err := s.repo.GetByID(uint(uid))
	if err != nil {
		return fmt.Errorf("用户不存在")
	}
	// 校验原密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword+user.Salt)); err != nil {
		return fmt.Errorf("原密码错误")
	}
	if oldPassword == newPassword {
		return fmt.Errorf("新密码不能与原密码相同")
	}
	salt := RandSalt()
	hash, err := models.HashPasswordWithSalt(newPassword, salt)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(uint(uid), hash, salt)
}

// GetDeptName 根据部门ID获取部门名称
func (s *UserServiceImpl) GetDeptName(deptID uint) (string, error) {
	if deptID <= 0 {
		return "", nil
	}
	var dept models.Dept
	if err := s.db.First(&dept, deptID).Error; err != nil {
		return "", err
	}
	return dept.Name, nil
}

// GetRoleNames 根据用户ID获取角色名称（逗号分隔）
func (s *UserServiceImpl) GetRoleNames(userID string) (string, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return "", err
	}
	if len(roles) == 0 {
		return "", nil
	}
	var names []string
	for _, r := range roles {
		names = append(names, r.Name)
	}
	return strings.Join(names, ","), nil
}

// UpdateUserProfile 修改个人中心用户信息
func (s *UserServiceImpl) UpdateUserProfile(userID string, nickname, avatar string, gender string, mobile, email string) error {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	updateMap := map[string]interface{}{
		"nickname": nickname,
		"avatar":   avatar,
		"gender":   gender,
		"mobile":   mobile,
		"email":    email,
	}
	return s.repo.UpdateUserProfile(uint(uid), updateMap)
}

func (s *UserServiceImpl) ListUserOptions(ctx *gin.Context) ([]models.UserOption, error) {
	return s.repo.ListUserOptions(ctx)
}

// UpdateLastLogin 更新用户最后登录信息
func (s *UserServiceImpl) UpdateLastLogin(userID uint, ClientIP string, loginTime time.Time, userAgent string) error {
	//解析userAgent
	device, os, browser := utils.ParseUserAgent(userAgent)
	return s.repo.UpdateLastLogin(userID, ClientIP, loginTime, userAgent, device, browser, os)
}
