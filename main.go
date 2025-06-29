package main

import (
	"golang-api/config"
	"golang-api/models"
	"golang-api/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func InitializeApp() *gin.Engine {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading env")
	}

	r := gin.Default()

	db := config.ConnectDatabase()

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Inventory{})
	db.AutoMigrate(&models.Order{})
	db.AutoMigrate(&models.OrderItem{})

	routes.SetupRoutes(r, db)

	return r
}

func main() {
	app := InitializeApp()
	app.Run(":8080")
}
