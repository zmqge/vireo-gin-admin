package routes

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// 自动扫描控制器方法中的Permission注解
func AutoBindPermissions(r *gin.Engine, controller interface{}) {
	t := reflect.TypeOf(controller)
	for i := 0; i < t.NumMethod(); i++ {
		_ = t.Method(i) // 明确忽略未使用的变量
		// 通过方法注释解析权限码（需配合go:generate）
		// 或使用结构体字段标签（推荐）
	}
}
