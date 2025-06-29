package routes

import (
	"golang-api/controllers"
	"golang-api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProductRoutes(router *gin.RouterGroup, db *gorm.DB) {
	productController := controllers.NewProductController(db)

	protected := router.Group("/")

	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/products", productController.CreateProduct)
		protected.GET("/products", productController.GetProduct)
		protected.GET("/products/:id", productController.GetProductByID)
		protected.PUT("/products/:id", productController.UpdateProduct)
		protected.DELETE("/products/:id", productController.DeleteProduct)

		protected.GET("/products/images/:fileName", controllers.DownloadFile)
	}
}
