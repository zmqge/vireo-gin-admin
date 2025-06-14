package middleware

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
	"gorm.io/gorm"
)

const (
	DataScopeAll     = 1 // 全部数据权限
	DataScopeDeptSub = 2 // 本部门及以下
	DataScopeDept    = 3 // 本部门
	DataScopeOwn     = 4 // 仅本人
	DataScopeCustom  = 5 // 自定义
)

func DATAPERM() gin.HandlerFunc {
	return func(c *gin.Context) {

		db, ok := c.MustGet("db").(*gorm.DB)
		if !ok {
			log.Printf("[DATAPERM] 数据库实例获取失败")
			response.Error(c, fmt.Errorf("内部服务错误"))
			c.Abort()
			return
		}

		// 1. 获取基础用户信息（不预加载关联）
		user, err := getBasicUser(c, db)
		if err != nil {
			log.Printf("[DATAPERM] 获取用户信息失败: %v", err)
			response.Error(c, fmt.Errorf("用户信息获取失败: %w", err))
			c.Abort()
			return
		}
		log.Printf("[DATAPERM] 基础用户信息: ID=%d, Name=%s, DeptID=%d",
			user.ID, user.Username, user.DeptID)
		log.Printf("[DATAPERM] 用户角色数量: %d", len(user.RoleList))
		for i, role := range user.RoleList {
			log.Printf("[DATAPERM] 角色[%d]: ID=%d, Name=%s, DataScope=%d",
				i, role.ID, role.Name, role.DataScope)
		}

		// 2. 按需加载数据
		if err := loadRequiredData(db, user); err != nil {
			log.Printf("[DATAPERM] 数据加载失败: %v", err)
			response.Error(c, fmt.Errorf("权限数据获取失败"))
			c.Abort()
			return
		}
		log.Printf("[DATAPERM] 用户 %s 数据权限加载成功, 范围: %d", user.Username, user.DataScope)

		// 3. 存储到上下文
		c.Set("currentUser", user)
		//用于判断是否使用了dataperm中间件
		c.Set("dataPermEnabled", true)

		// 验证存储结果
		if storedUser, exists := c.Get("currentUser"); exists {
			if u, ok := storedUser.(*models.User); ok {
				log.Printf("[DATAPERM] Context存储成功: ID=%d, Name=%s, DataScope=%d, Roles=%d",
					u.ID, u.Username, u.DataScope, len(u.RoleList))
			} else {
				log.Printf("[DATAPERM] Context类型错误: %T", storedUser)
			}
		} else {
			log.Printf("[DATAPERM] Context存储失败")
		}

		c.Next()
	}
}

// 获取基础用户信息（包含必要的关联数据）
func getBasicUser(c *gin.Context, db *gorm.DB) (*models.User, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("用户未登录")
	}

	var user models.User
	// 一次性加载用户、角色和部门信息
	if err := db.Preload("RoleList").Preload("Dept").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("用户查询失败: %w", err)
	}

	if len(user.RoleList) == 0 {
		return nil, fmt.Errorf("用户没有分配角色")
	}

	return &user, nil
}

// 按需加载必要数据
func loadRequiredData(db *gorm.DB, user *models.User) error {
	// 收集所有角色的权限范围
	deptMap := make(map[uint]struct{})
	hasPermissionSet := false
	maxScope := DataScopeOwn // 默认为仅本人数据

	log.Printf("[DATAPERM] 开始处理用户所有角色的数据权限...")

	for _, role := range user.RoleList {
		log.Printf("[DATAPERM] 处理角色权限: ID=%d, Name=%s, DataScope=%d",
			role.ID, role.Name, role.DataScope)

		if role.DataScope == DataScopeAll { // 如果有全部数据权限
			user.DataScope = DataScopeAll
			user.PermissionDepts = nil // 不需要限制部门
			log.Printf("[DATAPERM] 用户拥有全部数据权限，来自角色: %s", role.Name)
			return nil
		}

		if role.DataScope < maxScope {
			maxScope = role.DataScope
		}

		switch role.DataScope {
		case DataScopeDeptSub: // 本部门及以下
			if user.DeptID > 0 {
				var depts []models.Dept
				if err := db.Where("id = ? OR parent_id = ?", user.DeptID, user.DeptID).Find(&depts).Error; err != nil {
					log.Printf("[DATAPERM] 获取部门树失败: %v", err)
					continue
				}
				for _, d := range depts {
					deptMap[d.ID] = struct{}{}
				}
				hasPermissionSet = true
			}

		case DataScopeDept: // 本部门
			if user.DeptID > 0 {
				deptMap[user.DeptID] = struct{}{}
				hasPermissionSet = true
			}

		case DataScopeCustom: // 自定义部门
			deptIDs, err := role.GetCustomDepts()
			if err != nil {
				log.Printf("[DATAPERM] 获取自定义部门失败: %v", err)
				continue
			}
			for _, id := range deptIDs {
				deptMap[id] = struct{}{}
			}
			hasPermissionSet = true
		}
	}

	user.DataScope = maxScope

	// 如果没有任何部门权限，使用仅本人数据权限
	if !hasPermissionSet || len(deptMap) == 0 {
		user.DataScope = DataScopeOwn
		user.PermissionDepts = nil
		log.Printf("[DATAPERM] 用户没有有效的部门权限，使用仅本人数据权限")
		return nil
	}

	// 转换部门ID映射为切片
	var allDeptIDs []uint
	for id := range deptMap {
		allDeptIDs = append(allDeptIDs, id)
	}

	log.Printf("[DATAPERM] 用户最终可访问的部门ID: %v", allDeptIDs)
	user.PermissionDepts = allDeptIDs
	return nil
}
