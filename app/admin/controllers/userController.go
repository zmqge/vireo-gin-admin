package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/pkg/database"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
)

// UserController 用户控制器
// @Group(name="用户管理",path="/api/v1/")
type UserController struct {
	BaseController
	userService services.UserService
}

// RouteMeta 路由元数据

type RouteMeta struct {
	Method      string
	Path        string
	HandlerName string
	Middlewares []gin.HandlerFunc
}

// BaseController 控制器基类
type BaseController struct {
	Routes []RouteMeta
}

// NewUserController 创建 UserController 实例
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Route(method=GET, path="users/me", middlewares=["jwt"])
// @Permission(code="sys:user:me", name="获取当前用户信息", modules="个人中心", desc="获取当前登录用户的详细信息")
func (c *UserController) Me(ctx *gin.Context) {
	// 从上下文中获取用户 ID
	userID := ctx.GetString("userID") // 或者使用 ctx.Get("userID")

	// 调用服务层获取用户信息
	user, err := c.userService.GetUser(userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 获取用户的角色和权限
	permissionService := services.NewPermissionService()
	roles, permissions, err := permissionService.GetUserRolesAndPermissions(userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 返回用户信息
	response.Success(ctx, gin.H{
		"userId":   user.ID,       // 用户 ID
		"username": user.Username, // 用户名
		"nickname": user.Nickname, // 昵称
		"avatar":   user.Avatar,   // 头像
		"roles":    roles,         // 角色信息
		"perms":    permissions,   // 权限信息

	})

}

// @Route(method=DELETE, path="/users/:id", middlewares=["jwt","rbac"])
// @Permission(code="sys:user:delete",name="删除用户",modules="用户管理", desc="删除用户")
func (c *UserController) Delete(ctx *gin.Context) {
	if err := c.userService.Delete(ctx.Param("id")); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// @Route(method=GET, path="users/page",middlewares=["jwt","dataPerm"])
// @Permission(code="sys:user:page",name="用户分页列表", desc="获取用户分页列表",modules="用户管理")
// GetUserPage 获取用户分页列表
func (c *UserController) GetUserPage(ctx *gin.Context) {

	// 定义请求结构体，标记必填参数为required
	var req struct {
		PageNum  string `form:"pageNum" binding:"required"`  // 必填参数
		PageSize string `form:"pageSize" binding:"required"` // 必填参数

		// 可选参数
		Keywords   string   `form:"keywords"`
		Status     string   `form:"status"`
		RoleIDs    string   `form:"roleIds"`
		CreateTime []string `form:"createTime[]"` // 支持 createTime[0], createTime[1] 数组
		Field      string   `form:"field"`
		Direction  string   `form:"direction"`
		DeptID     string   `form:"deptId"`
	}
	// 绑定参数并验证必填项
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, fmt.Errorf("缺少必填参数: "+err.Error()))
		return
	}

	// 转换必填参数为int类型
	pageNum, err := strconv.Atoi(req.PageNum)
	if err != nil || pageNum <= 0 {
		response.Error(ctx, fmt.Errorf("pageNum必须是正整数"))
		return
	}

	pageSize, err := strconv.Atoi(req.PageSize)
	if err != nil || pageSize <= 0 {
		response.Error(ctx, fmt.Errorf("pageSize必须是正整数"))
		return
	}

	// 处理可选参数
	deptID := 0
	if req.DeptID != "" {
		if id, _ := strconv.Atoi(req.DeptID); id > 0 {
			deptID = id
		}
	}

	// 手动解析 createTime 数组参数，兼容 createTime[0]/createTime[1]
	createTimeStart := ctx.Query("createTime[0]")
	createTimeEnd := ctx.Query("createTime[1]")
	createTimeArr := []string{"", ""}
	if createTimeStart != "" || createTimeEnd != "" {
		createTimeArr[0] = createTimeStart
		createTimeArr[1] = createTimeEnd
	} else {
		createTimeArr = req.CreateTime
	}

	// 构建查询参数
	params := models.UserQueryParams{
		PageNum:    pageNum,
		PageSize:   pageSize,
		Keywords:   req.Keywords,
		Status:     req.Status,
		RoleIDs:    strings.Split(req.RoleIDs, ","),
		CreateTime: createTimeArr, // 优先用手动解析
		Field:      req.Field,
		Direction:  req.Direction,
		DeptID:     deptID,
	}

	// 调用服务层方法
	result, err := c.userService.GetUserPage(ctx, params)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 转换为响应格式
	userPageVOs := make([]models.UserPageVO, 0, len(result.Users))
	for _, user := range result.Users {
		// 查询部门名称
		deptName := ""
		if user.DeptID > 0 {
			var dept models.Dept
			db := database.GetDB()
			if db != nil && db.First(&dept, user.DeptID).Error == nil {
				deptName = dept.Name
			}
		}
		// 查询角色名称
		roleNames := ""
		roles, _ := c.userService.GetUserRoles(fmt.Sprintf("%d", user.ID))
		if len(roles) > 0 {
			var names []string
			for _, r := range roles {
				names = append(names, r.Name)
			}
			roleNames = strings.Join(names, ",")
		}
		// 格式化创建时间
		createTime := ""
		if !user.CreatedAt.IsZero() {
			createTime = user.CreatedAt.Format("2006-01-02 15:04:05")
		}
		userPageVOs = append(userPageVOs, models.UserPageVO{
			ID:         int64(user.ID),
			Username:   user.Username,
			Nickname:   user.Nickname,
			Mobile:     user.Mobile,
			Gender:     user.Gender,
			Avatar:     user.Avatar,
			Email:      user.Email,
			Status:     user.Status,
			DeptID:     user.DeptID,
			DeptName:   deptName,
			RoleNames:  roleNames,
			CreateTime: createTime,
		})
	}

	response.Success(ctx, models.DataUserPageVO{
		List:  userPageVOs,
		Total: result.Total,
	})
}

// 获取用户信息
// @Route(method=GET, path="/users/:id/form", middlewares=["jwt"])
// @Permission(code="sys:user:info",name="用户信息表单", modules="用户管理", desc="获取用户信息")
func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id") // 从URL参数中获取用户ID
	user, err := c.userService.GetUser(userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// 查询用户角色ID列表
	roles, err := c.userService.GetUserRoles(userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	var roleIds []int64
	for _, r := range roles {
		roleIds = append(roleIds, int64(r.ID))
	}
	resp := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"mobile":   user.Mobile,
		"gender":   user.Gender,
		"avatar":   user.Avatar,
		"email":    user.Email,
		"status":   user.Status,
		"deptId":   user.DeptID,
		"roleIds":  roleIds,
		"openId":   "",
	}

	response.Success(ctx, resp)
}

// 修改用户
// @Route(method=PUT, path="/users/:id", middlewares=["jwt"])
// @Permission(code="sys:user:edit",name="用户编辑", modules="用户管理", desc="更新用户信息")
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var req struct {
		Nickname string  `json:"nickname" binding:"required"`
		Mobile   string  `json:"mobile"`
		Gender   string  `json:"gender" binding:"required"`
		Avatar   string  `json:"avatar"`
		Email    string  `json:"email"`
		Status   int     `json:"status"`
		DeptID   int     `json:"deptId"`
		RoleIds  []int64 `json:"roleIds" binding:"required"`
		OpenId   string  `json:"openId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	// 调用 service 层
	if err := c.userService.UpdateUserFull(userID, req.Nickname, req.Mobile, req.Gender, req.Avatar, req.Email, req.Status, req.DeptID, req.RoleIds, req.OpenId); err != nil {
		response.Error(ctx, err)
		return
	}
	resp := map[string]interface{}{}
	response.Success(ctx, resp)
}

// 新增用户
// @Route(method=POST, path="/users", middlewares=["jwt"])
// @Permission(code="sys:user:add",name="用户新增", modules="用户管理", desc="创建新用户")
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req struct {
		Username string  `json:"username" binding:"required"`
		Nickname string  `json:"nickname" binding:"required"`
		Mobile   string  `json:"mobile"`
		Gender   string  `json:"gender" binding:"required"`
		Avatar   string  `json:"avatar"`
		Email    string  `json:"email"`
		Status   int     `json:"status"`
		DeptID   uint    `json:"deptId"`
		RoleIds  []int64 `json:"roleIds" binding:"required"`
		OpenId   string  `json:"openId"`
		Password string  `json:"password"` // 可选，若有初始密码
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Password == "" {
		req.Password = "123456" // 默认初始密码
	}
	// 调用 service 层
	err := c.userService.CreateUserFull(req.Username, req.Nickname, req.Mobile, req.Gender, req.Avatar, req.Email, req.Status, req.DeptID, req.RoleIds, req.OpenId, req.Password)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	resp := map[string]interface{}{}
	response.Success(ctx, resp)
}

// 重置用户密码
// @Route(method=PUT, path="/users/:id/password/reset", middlewares=["jwt"])
// @Permission(code="sys:user:reset-password",name="重置密码",modules="用户管理", desc="重置用户密码")
func (c *UserController) ResetPassword(ctx *gin.Context) {
	userID := ctx.Param("id")
	password := ctx.Query("password")
	if userID == "" || password == "" {
		response.BadRequest(ctx, "userId 和 password 不能为空")
		return
	}
	err := c.userService.ResetPassword(userID, password)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	resp := map[string]interface{}{}
	response.Success(ctx, resp)
}

// 获取个人中心用户信息
// @Route(method=GET, path="/users/profile", middlewares=["jwt"])
// @Permission(code="sys:user:profile", name="个人信息", modules="个人中心", desc="获取当前登录用户的个人中心信息")
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		response.Unauthorized(ctx, "用户未认证")
		return
	}
	user, err := c.userService.GetUser(userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// 通过 service 获取部门名称
	deptName, _ := c.userService.GetDeptName(user.DeptID)
	// 通过 service 获取角色名称
	roleNames, _ := c.userService.GetRoleNames(userID)
	createTime := ""
	if !user.CreatedAt.IsZero() {
		createTime = user.CreatedAt.Format("2006-01-02 15:04:05")
	}
	resp := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"gender":     user.Gender,
		"mobile":     user.Mobile,
		"email":      user.Email,
		"deptName":   deptName,
		"roleNames":  roleNames,
		"createTime": createTime,
	}
	response.Success(ctx, resp)
}

// 修改当前用户密码
// @Route(method=PUT, path="/users/password", middlewares=["jwt"])
// @Permission(code="sys:user:change-password", name="修改密码", modules="个人中心", desc="修改当前登录用户的密码")
// ChangePassword 修改当前登录用户的密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		response.Unauthorized(ctx, "用户未认证")
		return
	}
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	if req.OldPassword == req.NewPassword {
		response.BadRequest(ctx, "新密码不能与原密码相同")
		return
	}
	if err := c.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// 修改当前登录用户的个人中心信息
// @Route(method=PUT, path="/users/profile", middlewares=["jwt"])
// @Permission(code="sys:user:update-profile", name="修改个人信息", modules="个人中心", desc="修改当前登录用户的个人中心信息")
func (c *UserController) UpdateMyProfile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		response.Unauthorized(ctx, "用户未认证")
		return
	}
	var req struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Gender   string `json:"gender"`
		Mobile   string `json:"mobile"`
		Email    string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}
	// 只允许本人操作自己的信息
	if req.ID > 0 && strconv.FormatInt(req.ID, 10) != userID {
		response.Forbidden(ctx, "无权修改他人信息")
		return
	}
	err := c.userService.UpdateUserProfile(userID, req.Nickname, req.Avatar, req.Gender, req.Mobile, req.Email)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil)
}

// @Summary 获取用户下拉选项
// @Tags 用户管理
// @Route(method=GET, path="/users/options", middlewares=["jwt","dataPerm"])
// @Permission(code="sys:user:options", name="用户下拉选项", modules="用户管理", desc="获取用户下拉选项")
func (c *UserController) ListUserOptions(ctx *gin.Context) {
	options, err := c.userService.ListUserOptions(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	if len(options) == 0 {
		response.Success(ctx, []interface{}{})
		return
	}
	response.Success(ctx, options)
}
