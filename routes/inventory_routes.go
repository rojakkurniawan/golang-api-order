package routes

import (
	"golang-api/controllers"
	"golang-api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupInventoryRoutes(router *gin.RouterGroup, db *gorm.DB) {
	inventoryController := controllers.NewInventoryController(db)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/inventory", inventoryController.CreateInventory)
		protected.GET("/inventory", inventoryController.GetInventories)
		protected.GET("/inventory/:id", inventoryController.GetInventoryByID)
		protected.PUT("/inventory/:id", inventoryController.UpdateInventory)
		protected.DELETE("/inventory/:id", inventoryController.DeleteInventory)

		protected.PUT("/inventory/stock", inventoryController.UpdateStock)
	}
}
