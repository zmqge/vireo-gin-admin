package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

)

func RegisterAllRoutes(engine *gin.Engine, db *gorm.DB) {
	RegisterAdminRoutes(engine, db)
}
