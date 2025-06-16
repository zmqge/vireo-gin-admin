package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"github.com/zmqge/vireo-gin-admin/config"
	"github.com/zmqge/vireo-gin-admin/pkg/auth"
	"github.com/zmqge/vireo-gin-admin/pkg/redis"
	"github.com/zmqge/vireo-gin-admin/pkg/response"
	"github.com/zmqge/vireo-gin-admin/utils"
)

// @Group(path="/api/v1/auth", name="认证", desc="认证接口")
type AuthController struct {
	userService  services.UserService
	tokenService services.TokenService
}

func NewAuthController(
	userService services.UserService,
	tokenService services.TokenService,
) *AuthController {
	return &AuthController{
		userService:  userService,
		tokenService: tokenService,
	}
}

// @Route(path="/login", method="POST", desc="登录")
func (c *AuthController) Login(ctx *gin.Context) {
	// 1. 定义请求参数结构
	type loginRequest struct {
		Username   string `json:"username" form:"username" binding:"required"`
		Password   string `json:"password" form:"password" binding:"required"`
		CaptchaID  string `json:"captchaKey" form:"captchaKey" binding:"required"`
		CaptchaAns string `json:"captchaCode" form:"captchaCode" binding:"required"`
	}

	// 2. 绑定参数
	var req loginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest(ctx, "无效的请求参数")
		return
	}

	// 3. 验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaAns) {
		response.BadRequest(ctx, "验证码错误")
		return
	}

	// 4. 验证用户名和密码
	user, err := c.userService.VerifyUser(req.Username, req.Password)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 5. 生成双Token
	jwt := auth.NewJWT()
	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Username)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	lastLoginTime := time.Now()
	clientIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	// 6. 存储Refresh Token到数据库（可选）
	if err := c.tokenService.SaveRefreshToken(user.ID, refreshToken, clientIP, lastLoginTime); err != nil {
		response.Error(ctx, err)
		return
	}

	// 更新用户最后登录信息
	if err := c.userService.UpdateLastLogin(user.ID, clientIP, lastLoginTime, userAgent); err != nil {
		response.Error(ctx, err)
		log.Println("更新用户最后登录信息失败:", err)
	}

	// 7. 返回双Token
	response.Success(ctx, gin.H{
		"tokenType":    "Bearer",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expiresIn":    config.App.JWT.AccessExpire.Seconds(),
	})
}

// @Route(method=GET, path="/captcha", middlewares=[])
func (c *AuthController) GetCaptcha(ctx *gin.Context) {
	// 生成验证码
	id, b64s, err := utils.GenerateCaptcha()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": "99999", // 自定义错误码
			"msg":  "生成验证码失败",
		})
		return
	}

	// 返回验证码 ID 和图片
	response.Success(ctx, gin.H{
		"captchaKey":    id,
		"captchaBase64": b64s,
	}, "一切ok")
}

// 退出登陆
// @Route(method=DELETE, path="/logout", middlewares=["jwt"])
func (c *AuthController) Logout(ctx *gin.Context) {
	// 从请求头或参数中获取Refresh Token
	userID := ctx.GetString("userID")       // 从上下文中获取用户ID
	userIDUint, err := strconv.Atoi(userID) // 转换为uint类型
	if err != nil {
		response.BadRequest(ctx, "无效的用户ID")
		return
	}
	// 从Redis中删除Refresh Token
	err = c.tokenService.DeleteRefreshToken(uint(userIDUint))
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 返回成功响应
	response.Success(ctx, nil, "退出成功")

}

func getUserIDFromRefreshToken(refreshToken string) (uint, error) {
	// 从Redis查询refreshToken对应的userID
	keys, err := redis.Client.Keys(context.Background(), "refresh:*").Result()
	if err != nil {
		return 0, err
	}

	for _, key := range keys {
		storedToken, err := redis.Client.Get(context.Background(), key).Result()
		if err != nil {
			continue
		}
		if storedToken == refreshToken {
			// 从key中提取userID
			userIDStr := strings.TrimPrefix(key, "refresh:")
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				return 0, err
			}
			return uint(userID), nil
		}
	}
	return 0, errors.New("refresh token not found")
}

// @Route(method=POST, path="/refresh-token", middlewares=[])
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("RefreshToken")
	if refreshToken == "" {
		refreshToken = ctx.Query("refreshToken")
	}
	if refreshToken == "" {
		response.BadRequest(ctx, "Refresh Token缺失")
		return
	}

	// 从Redis验证Refresh Token
	userID, err := getUserIDFromRefreshToken(refreshToken)
	if err != nil {
		response.RefresTokenExpired(ctx, "Refresh Token无效或已过期")
		return
	}

	// 生成新的 Access Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.App.JWT.AccessExpire)),
	})

	// 签名并生成 Token 字符串
	newAccessToken, err := token.SignedString([]byte(config.App.JWT.AccessSecret))
	if err != nil {
		response.Error(ctx, errors.New("生成 Access Token 失败"))
		return
	}

	// 返回新的 Access Token和调试信息
	response.Success(ctx, gin.H{
		"accessToken":  newAccessToken,
		"refreshToken": refreshToken,
		"tokenType":    "Bearer",
		"expiresIn":    config.App.JWT.AccessExpire.Seconds(),
	})
}
