package scopes

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/config"
	"gorm.io/gorm"
)

// DataPermissionScope 数据权限过滤器
func DataPermissionScope(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 有使用数据权限过滤器，显示所有数据
		if !ctx.GetBool("dataPermEnabled") {
			log.Printf("[DataPermissionScope]:没有使用数据权限过滤器 ")
			return db // 没有使用数据权限过滤器，显示所有数据
		}

		// 1. 获取当前用户
		user, err := getCurrentUser(ctx)
		if err != nil {
			log.Printf("[DataPermissionScope] 获取用户失败: %v", err)
			return db.Where("1 = 0") // 无权限
		}

		// 2. 管理员检查或全部数据权限
		if isAdmin(user) || user.DataScope == 1 { // 1 = 全部数据权限
			log.Printf("[DataPermissionScope] 用户 %s 拥有全部数据权限", user.Username)
			return db
		}

		// 3. 检查请求参数中的deptId，如果有，则检查权限，属于可访问权限放行，如果没有部门权限，那就是只能访问自己的数据，那么控制自己的数据权限。
		deptIdStr := ctx.Query("deptId")
		if deptIdStr != "" {
			deptId, err := strconv.ParseUint(deptIdStr, 10, 64)
			if err == nil {
				// 检查deptId是否在用户权限部门列表中
				for _, permDept := range user.PermissionDepts {
					if uint(deptId) == permDept {
						log.Printf("[DataPermissionScope] 用户 %s 通过参数deptId访问部门 %d", user.Username, deptId)
						return db
					}
				}
				return db.Where("creator_id = ?", user.ID)

			}
		}

		// 4. 检查部门权限列表，自己的数据是最小权限
		// 如果没有部门权限列表，则只能查看自己的数据
		if len(user.PermissionDepts) == 0 {
			log.Printf("[DataPermissionScope] 用户 %s 只能查看本人数据", user.Username)
			return db.Where("creator_id = ?", user.ID)
		}

		// 有部门权限列表，则可以查看指定部门的数据
		log.Printf("[DataPermissionScope] 用户 %s 可访问部门: %v", user.Username, user.PermissionDepts)
		return db.Where("dept_id IN ?", user.PermissionDepts)
	}
}

// getCurrentUser 从上下文中获取用户信息
func getCurrentUser(ctx *gin.Context) (*models.User, error) {
	user, exists := ctx.Get("currentUser")
	if !exists {
		return nil, errors.New("用户信息未找到")
	}

	currentUser, ok := user.(*models.User)
	if !ok || currentUser == nil {
		return nil, errors.New("用户信息类型错误")
	}

	if len(currentUser.RoleList) == 0 {
		return nil, errors.New("用户角色信息未加载")
	}

	return currentUser, nil
}

// isAdmin 检查是否是超级管理员（支持多角色）
func isAdmin(user *models.User) bool {
	superAdminRole := config.App.RBAC.SuperAdminRole
	if superAdminRole == "" {
		superAdminRole = "super_admin" // 默认超级管理员角色名
	}
	for _, r := range user.RoleList {
		if r.Name == superAdminRole {
			return true
		}
	}
	return false
}
