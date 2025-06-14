package middleware

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
	"gorm.io/gorm"
)

// RBAC 创建一个基于角色的访问控制中间件
// requiredCodes: 需要的权限码列表
func RBAC(requiredCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取数据库实例
		db, err := getDBFromContext(c)
		if err != nil {
			log.Println("数据库实例获取失败:", err)
			response.Error(c, fmt.Errorf("内部服务错误"))
			c.Abort()
			return
		}

		// 2. 获取并验证用户ID
		userID, err := getUserIDFromContext(c)
		if err != nil {
			log.Println("用户ID获取失败:", err)
			response.Unauthorized(c, "请先登录或登录信息无效")
			c.Abort()
			return
		}

		// 3. 检查是否是超级管理员
		if isSuperAdmin(db, userID) {
			c.Next()
			return
		}

		// 4. 权限检查（如果需要特定权限）
		if len(requiredCodes) > 0 {
			if err := checkUserPermissions(db, userID, requiredCodes); err != nil {
				log.Println("权限检查失败:", err)
				switch err {
				case ErrPermissionDenied:
					response.Forbidden(c, "无权访问")
				default:
					response.Error(c, fmt.Errorf("权限验证失败"))
				}
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// 定义错误类型
var (
	ErrInvalidDBInstance = fmt.Errorf("无效的数据库实例")
	ErrInvalidUserID     = fmt.Errorf("无效的用户ID")
	ErrPermissionDenied  = fmt.Errorf("权限不足")
)

// getDBFromContext 从上下文中获取数据库实例
func getDBFromContext(c *gin.Context) (*gorm.DB, error) {
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		return nil, ErrInvalidDBInstance
	}
	return db, nil
}

// getUserIDFromContext 从上下文中获取并解析用户ID
func getUserIDFromContext(c *gin.Context) (uint, error) {
	userIDVal := c.GetString("userID")
	if userIDVal == "" {
		return 0, fmt.Errorf("用户未登录")
	}
	return parseUserID(userIDVal)
}

// parseUserID 解析用户ID为uint类型
func parseUserID(userIDVal interface{}) (uint, error) {
	switch v := userIDVal.(type) {
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case float64:
		return uint(v), nil
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%w: %v", ErrInvalidUserID, err)
		}
		return uint(n), nil
	default:
		return 0, fmt.Errorf("%w: 不支持的类型 %T", ErrInvalidUserID, v)
	}
}

// checkUserPermissions 检查用户是否拥有所有需要的权限
func checkUserPermissions(db *gorm.DB, userID uint, requiredCodes []string) error {
	perms, err := getUserPermissions(db, userID)
	if err != nil {
		return fmt.Errorf("获取用户权限失败: %w", err)
	}

	if !hasAllPermissions(perms, requiredCodes) {
		return fmt.Errorf("%w: 需要 %v, 实际 %v", ErrPermissionDenied, requiredCodes, perms)
	}
	return nil
}

// getUserPermissions 获取用户的所有权限码
func getUserPermissions(db *gorm.DB, userID uint) ([]string, error) {
	var perms []string
	err := db.Table("user_roles").
		Select("DISTINCT role_permissions.permission_code").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permission_code", &perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}

// hasAllPermissions 检查用户是否拥有所有需要的权限
func hasAllPermissions(userPerms, requiredPerms []string) bool {
	permSet := make(map[string]struct{}, len(userPerms))
	for _, perm := range userPerms {
		permSet[perm] = struct{}{}
	}

	for _, reqPerm := range requiredPerms {
		if _, ok := permSet[reqPerm]; !ok {
			return false
		}
	}
	return true
}

// isSuperAdmin 检查用户是否是超级管理员
func isSuperAdmin(db *gorm.DB, userID uint) bool {
	roles, err := getUserRoleNames(db, userID)
	if err != nil {
		log.Printf("获取用户角色失败: %v", err)
		return false
	}

	superAdminRole := config.App.RBAC.SuperAdminRole
	if superAdminRole == "" {
		superAdminRole = "super_admin" // 默认超级管理员角色名
	}

	for _, role := range roles {
		if role == superAdminRole {
			return true
		}
	}
	return false
}

// getUserRoleNames 获取用户的所有角色名称
func getUserRoleNames(db *gorm.DB, userID uint) ([]string, error) {
	var roles []string
	err := db.Table("user_roles").
		Select("DISTINCT roles.name").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("name", &roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
