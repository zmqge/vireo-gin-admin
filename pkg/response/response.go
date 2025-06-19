package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data interface{}, msg ...string) {
	message := "success"
	if len(msg) > 0 {
		message = msg[0]
	}
	c.JSON(200, gin.H{"code": 0, "data": data, "msg": message})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(403, gin.H{"code": 403, "msg": msg})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(400, gin.H{"code": 400, "msg": msg})
}

func Error(c *gin.Context, err error) {
	c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
}

func PageSuccess(c *gin.Context, data interface{}, total int64) {
	c.JSON(200, gin.H{"code": 0, "data": data, "total": total})
}

// Unauthorized 返回401未授权错误
func Unauthorized(ctx *gin.Context, message string) {
	ctx.JSON(401, gin.H{
		"code": 401,
		"msg":  message,
	})
}

// 刷新token过期，返回402错误
func RefresTokenExpired(c *gin.Context, msg string) {
	c.JSON(402, gin.H{"code": 402, "msg": msg})
}

// NotFound 返回404未找到错误
func NotFound(c *gin.Context, msg string) {
	c.JSON(404, gin.H{"code": 404, "msg": msg})
}
func DemoMode(c *gin.Context, msg string) {
	c.JSON(403, gin.H{"code": 403, "msg": msg})
}
