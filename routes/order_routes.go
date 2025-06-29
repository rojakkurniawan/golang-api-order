package routes

import (
	"golang-api/controllers"
	"golang-api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRoutes(router *gin.RouterGroup, db *gorm.DB) {
	orderController := controllers.NewOrderController(db)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/orders", orderController.CreateOrder)
		protected.GET("/orders", orderController.GetOrders)
		protected.GET("/orders/:id", orderController.GetOrderByID)
		protected.PUT("/orders/:id/status", orderController.UpdateOrderStatus)
		protected.DELETE("/orders/:id", orderController.DeleteOrder)
	}
}
