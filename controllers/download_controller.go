package controllers

import (
	"golang-api/models"
	"golang-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DownloadFile(c *gin.Context) {
	fileName := c.Param("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "file name is required",
		})
		return
	}

	filePath, err := services.DownloadFile(fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.File(filePath)
}
