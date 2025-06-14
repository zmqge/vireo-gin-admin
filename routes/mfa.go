package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zmqge/vireo-gin-admin/app/admin/controllers"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
	"gorm.io/gorm"
)

// 确保函数名首字母大写（导出）
func RegisterMFARoutes(r *gin.Engine, db *gorm.DB) {
	mfaService := services.NewMFA(db)
	mfaController := controllers.NewMFA(mfaService)

	mfaGroup := r.Group("/mfa")
	{
		mfaGroup.POST("/enable", mfaController.Enable)
	}

	routes := r.Routes()
	for _, route := range routes {
		logrus.Info("已注册路由",
			logrus.Fields{
				"path":   route.Path,
				"method": route.Method,
			},
		)
	}
}
