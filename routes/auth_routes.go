package routes

import (
	"golang-api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(router *gin.RouterGroup, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	protected := router.Group("/")
	{
		protected.POST("/register", authController.Register)
		protected.POST("/login", authController.Login)
	}

}
