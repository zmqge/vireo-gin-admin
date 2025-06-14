Gin + RBAC 后台管理系统完整实现方案 ,系统名称为：vireo-gin-admin，系统采用前后端分离架构，后端采用Gin框架，前端采用Vue3 + Vite + Element Plus。
一、系统架构设计
1.1 技术栈选型
- 后端框架: Gin (轻量级高性能HTTP框架)
- 数据库: MySQL 8.0+ (支持JSON字段和窗口函数)
- 缓存: Redis (用于JWT黑名单和权限缓存)
- ORM: GORM (支持事务和复杂查询)
- 认证: JWT + 双Token机制(access/refresh)
- API文档: Swagger UI (通过注释自动生成)
1.2 分层架构
markdown
复制
.
├── app
│   ├── controllers      # 控制器层
│   ├── middleware       # 中间件层
│   ├── models           # 数据模型层
│   ├── services         # 业务逻辑层
│   └── validators       # 请求验证层
├── config               # 配置管理
├── database             # 数据库脚本
├── docs                 # Swagger文档
├── pkg
│   ├── auth             # 认证模块
│   ├── constant         # 全局常量
│   ├── errors           # 错误处理
│   ├── logger           # 日志模块
│   └── response         # 统一响应
├── routes               # 路由定义
├── scripts              # 部署脚本
├── storage              # 文件存储
└── utils                # 工具函数
二、核心数据库设计
2.1 完整ER图
mermaid
复制
erDiagram
    users ||--o{ user_roles : "多对多"
    users {
        bigint id PK
        varchar(50) username
        varchar(255) password
        varchar(50) salt
        tinyint status
        datetime created_at
        datetime updated_at
    }
    
    roles ||--o{ user_roles : "多对多"
    roles ||--o{ role_permissions : "一对多"
    roles {
        bigint id PK
        varchar(50) name
        varchar(255) desc
        datetime created_at
    }
    
    permissions ||--o{ role_permissions : "一对多"
    permissions ||--o{ menu_permissions : "一对多"
    permissions {
        bigint id PK
        varchar(50) code
        varchar(50) name
        enum('api','menu','button') type
        varchar(255) desc
    }
    
    menus ||--o{ menu_permissions : "一对多"
    menus {
        bigint id PK
        bigint parent_id
        varchar(50) name
        varchar(100) path
        varchar(100) component
        varchar(50) icon
        int sort
        bool hidden
    }
    
    user_roles {
        bigint user_id PK,FK
        bigint role_id PK,FK
    }
    
    role_permissions {
        bigint role_id PK,FK
        bigint permission_id PK,FK
    }
    
    menu_permissions {
        bigint menu_id PK,FK
        bigint permission_id PK,FK
    }
2.2 初始化SQL脚本
sql
复制
-- 系统初始化数据
INSERT INTO permissions (code, name, type, `desc`) VALUES 
('system:admin', '超级管理员权限', 'api', '拥有所有权限'),
('user:manage', '用户管理权限', 'menu', '用户管理菜单入口'),
('user:create', '创建用户', 'api', '创建用户权限'),
('user:update', '编辑用户', 'api', '编辑用户权限'),
('user:delete', '删除用户', 'api', '删除用户权限');
INSERT INTO roles (name, `desc`) VALUES 
('super_admin', '超级管理员'),
('normal_user', '普通用户');
-- 超级管理员拥有所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;
-- 默认管理员用户
INSERT INTO users (username, password, salt, status) VALUES 
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrqop6wWZklsZBSjT5Z.0S3JQ8b5.6O', 'abcdef', 1);
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);
三、核心模块实现
3.1 JWT认证模块 (pkg/auth/jwt.go)
go
复制
package auth
import (
        "time"
        "github.com/golang-jwt/jwt/v5"
)
type JWT struct {
        SigningKey []byte
        ExpireTime time.Duration
        BufferTime time.Duration
}
type CustomClaims struct {
        UserID   uint
        Username string
        Roles    []string
        jwt.RegisteredClaims
}
func NewJWT(signKey string, expire, buffer time.Duration) *JWT {
        return &JWT{
                SigningKey: []byte(signKey),
                ExpireTime: expire,
                BufferTime: buffer,
        }
}
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return token.SignedString(j.SigningKey)
}
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
        token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
                return j.SigningKey, nil
        })
        
        if err != nil {
                return nil, err
        }
        
        if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
                return claims, nil
        }
        
        return nil, err
}
// 双Token生成
func (j *JWT) GenerateTokens(user *models.User) (accessToken, refreshToken string, err error) {
        claims := CustomClaims{
                UserID:   user.ID,
                Username: user.Username,
                Roles:    user.GetRoleCodes(),
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
                        IssuedAt:  jwt.NewNumericDate(time.Now()),
                        Issuer:    "gin-admin",
                },
        }
        
        accessToken, err = j.CreateToken(claims)
        if err != nil {
                return "", "", err
        }
        
        // Refresh token有更长有效期
        claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.ExpireTime * 2))
        refreshToken, err = j.CreateToken(claims)
        
        return
}
3.2 RBAC中间件 (middleware/rbac.go)
go
复制
package middleware
import (
        "net/http"
        "strings"
        
        "gin-admin/pkg/auth"
        "gin-admin/pkg/response"
        
        "github.com/gin-gonic/gin"
)
// RBAC 权限检查中间件
func RBAC(permissionCode string) gin.HandlerFunc {
        return func(c *gin.Context) {
                claims, exists := c.Get("claims")
                if !exists {
                        response.Unauthorized(c, "认证信息不存在")
                        c.Abort()
                        return
                }
                
                // 超级管理员跳过检查
                if auth.HasRole(claims.(*auth.CustomClaims).Roles, "super_admin") {
                        c.Next()
                        return
                }
                
                // 检查权限
                if !auth.CheckPermission(c.Request.Method, c.FullPath(), permissionCode) {
                        response.Forbidden(c, "没有访问权限")
                        c.Abort()
                        return
                }
                
                c.Next()
        }
}
// 权限缓存检查
func CheckPermission(method, path, code string) bool {
        // 1. 从Redis缓存检查
        cacheKey := fmt.Sprintf("perm:%s:%s:%s", method, path, code)
        if redis.Exists(cacheKey) {
                return redis.GetBool(cacheKey)
        }
        
        // 2. 数据库检查
        var count int64
        db.Model(&models.RolePermission{}).
                Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
                Where("permissions.code = ? AND role_permissions.role_id IN (?)", 
                        code, getUserRoleIDs()).
                Count(&count)
        
        // 缓存结果
        redis.Set(cacheKey, count > 0, 30*time.Minute)
        return count > 0
}
3.3 权限服务层 (services/permission.go)
go
复制
package services
import (
        "gin-admin/app/models"
        "gin-admin/pkg/tree"
)
type PermissionService struct {
        DB *gorm.DB
}
// 获取用户权限树
func (s *PermissionService) GetUserPermissionTree(userID uint) (*tree.Node, error) {
        var permissions []models.Permission
        err := s.DB.Model(&models.User{}).
                Joins("JOIN user_roles ON user_roles.user_id = users.id").
                Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
                Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
                Where("users.id = ?", userID).
                Find(&permissions).Error
        
        if err != nil {
                return nil, err
        }
        
        // 构建权限树
        root := tree.NewNode("root", "根节点")
        for _, perm := range permissions {
                node := &tree.Node{
                        ID:    perm.Code,
                        Name:  perm.Name,
                        Type:  perm.Type,
                        Meta:  perm,
                }
                root.AddChild(node)
        }
        
        return root, nil
}
// 检查API权限
func (s *PermissionService) CheckAPIPermission(userID uint, method, path string) bool {
        var count int64
        s.DB.Model(&models.User{}).
                Joins("JOIN user_roles ON user_roles.user_id = users.id").
                Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
                Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
                Where("users.id = ? AND permissions.type = 'api' AND permissions.code = ?", 
                        userID, generatePermissionCode(method, path)).
                Count(&count)
        
        return count > 0
}
// 生成权限码 (GET:/api/v1/users -> get:api:v1:users)
func generatePermissionCode(method, path string) string {
        path = strings.TrimPrefix(path, "/")
        path = strings.ReplaceAll(path, "/", ":")
        return strings.ToLower(method) + ":" + path
}
3.4 用户控制器 (controllers/user.go)
go
复制
package controllers
import (
        "gin-admin/app/services"
        "gin-admin/pkg/response"
        
        "github.com/gin-gonic/gin"
)
type UserController struct {
        userService *services.UserService
}
// @Summary 获取用户列表
// @Tags 用户管理
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} response.PageResult
// @Router /users [get]
// @Permission user:list
func (c *UserController) List(ctx *gin.Context) {
        var req struct {
                Page int `form:"page"`
                Size int `form:"size"`
        }
        
        if err := ctx.ShouldBindQuery(&req); err != nil {
                response.BadRequest(ctx, err.Error())
                return
        }
        
        users, total, err := c.userService.GetUserList(req.Page, req.Size)
        if err != nil {
                response.Error(ctx, err)
                return
        }
        
        response.PageSuccess(ctx, users, total)
}
// @Summary 创建用户
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Param data body services.CreateUserRequest true "用户信息"
// @Success 200 {object} response.Result
// @Router /users [post]
// @Permission user:create
func (c *UserController) Create(ctx *gin.Context) {
        var req services.CreateUserRequest
        if err := ctx.ShouldBindJSON(&req); err != nil {
                response.BadRequest(ctx, err.Error())
                return
        }
        
        if err := c.userService.CreateUser(&req); err != nil {
                response.Error(ctx, err)
                return
        }
        
        response.Success(ctx)
}
四、路由配置 (routes/admin.go)
go
复制
package routes
import (
        "gin-admin/app/controllers"
        "gin-admin/middleware"
        
        "github.com/gin-gonic/gin"
)
func RegisterAdminRoutes(r *gin.Engine, ctrl *controllers.AdminController) {
        // 认证路由组
        auth := r.Group("/admin")
        auth.Use(middleware.JWTAuth())
        {
                // 用户管理
                user := auth.Group("/users")
                user.Use(middleware.RBAC("user:manage"))
                {
                        user.GET("", ctrl.User.List)       // 需要 user:list 权限
                        user.POST("", ctrl.User.Create)    // 需要 user:create 权限
                        user.PUT("/:id", ctrl.User.Update) // 需要 user:update 权限
                }
                
                // 角色管理
                role := auth.Group("/roles")
                role.Use(middleware.RBAC("role:manage"))
                {
                        role.GET("", ctrl.Role.List)
                        role.POST("/:id/permissions", ctrl.Role.AssignPermissions)
                }
                
                // 权限管理
                perm := auth.Group("/permissions")
                perm.Use(middleware.RBAC("perm:manage"))
                {
                        perm.GET("/tree", ctrl.Permission.GetTree)
                }
        }
        
        // 公开路由
        public := r.Group("/admin")
        {
                public.POST("/login", ctrl.Auth.Login)
                public.POST("/refresh", ctrl.Auth.RefreshToken)
        }
}
五、前端集成方案
5.1 权限指令 (Vue示例)
javascript
复制
// directives/permission.js
export default {
  inserted(el, binding, vnode) {
    const { value } = binding
    const permissions = store.getters.permissions
    
    if (value && value instanceof Array && value.length > 0) {
      const hasPermission = permissions.some(perm => {
        return value.includes(perm.code)
      })
      
      if (!hasPermission) {
        el.parentNode && el.parentNode.removeChild(el)
      }
    } else {
      throw new Error(`需要权限代码! 例如 v-permission="['user:create']"`)
    }
  }
}
5.2 动态路由生成
javascript
复制
// router/index.js
function filterAsyncRoutes(routes, permissions) {
  const res = []
  
  routes.forEach(route => {
    const tmp = { ...route }
    if (hasPermission(permissions, tmp)) {
      if (tmp.children) {
        tmp.children = filterAsyncRoutes(tmp.children, permissions)
      }
      res.push(tmp)
    }
  })
  
  return res
}
function hasPermission(permissions, route) {
  if (route.meta && route.meta.permission) {
    return permissions.some(perm => perm.code === route.meta.permission)
  }
  return true
}

七、扩展功能建议
1. 数据权限：基于部门的数据过滤
2. go
3. 复制
func (s *UserService) GetUserList(deptID uint) ([]User, error) {
    if !auth.CheckDataPermission(deptID) {
        return nil, errors.New("无数据权限")
    }
    // 查询逻辑...
}
1. 操作日志：记录关键操作
2. go
3. 复制
type OperationLog struct {
    UserID    uint
    IP        string
    Method    string
    Path      string
    Status    int
    UserAgent string
}
1. 权限模板：预定义角色权限组
2. sql
3. 复制
CREATE TABLE permission_templates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50),
    permissions JSON COMMENT '权限ID数组'
);
1. 权限回收：实时撤销权限
2. go
3. 复制
func RevokePermission(userID uint) error {
    // 清除Redis权限缓存
    redis.Delete(fmt.Sprintf("user:%d:permissions", userID))
    // 标记JWT失效
    redis.Set(fmt.Sprintf("jwt:invalid:%d", userID), true, 24*time.Hour)
}
这个完整方案提供了从数据库设计到前端集成的全链路实现，您可以根据实际需求进行裁剪或扩展。关键特点是：
1. 权限与菜单解耦，灵活支持多种场景
2. 采用双Token机制增强安全性
3. 完善的缓存策略提高性能
4. 清晰的代码分层和模块划分
5. 完整的API文档和部署方案