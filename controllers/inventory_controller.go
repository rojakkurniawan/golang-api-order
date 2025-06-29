package controllers

import (
	"fmt"
	"golang-api/models"
	"golang-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryController struct {
	InventoryService *services.InventoryService
}

func NewInventoryController(db *gorm.DB) *InventoryController {
	return &InventoryController{
		InventoryService: services.NewInventoryService(db),
	}
}

func (ic *InventoryController) CreateInventory(c *gin.Context) {
	var req models.AddInventoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid data: " + err.Error(),
		})
		return
	}

	newInventory := models.Inventory{
		ProductID: req.ProductID,
		Jumlah:    req.Jumlah,
		Lokasi:    req.Lokasi,
	}

	inventory, err := ic.InventoryService.CreateInventory(&newInventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := models.InventoryResponse{
		ID:        inventory.ID,
		ProductID: inventory.ProductID,
		Product: models.ProductResponse{
			ID:         inventory.Product.ID,
			Nama:       inventory.Product.Nama,
			Deskripsi:  inventory.Product.Deskripsi,
			Harga:      inventory.Product.Harga,
			Kategori:   inventory.Product.Kategori,
			FotoProduk: inventory.Product.FotoProduk,
			CreatedAt:  inventory.Product.CreatedAt,
			UpdatedAt:  inventory.Product.UpdatedAt,
		},
		Jumlah:    inventory.Jumlah,
		Lokasi:    inventory.Lokasi,
		CreatedAt: inventory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: inventory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Inventory successfully created",
		Data:    response,
	})
}

func (ic *InventoryController) GetInventories(c *gin.Context) {
	var req models.GetInventoryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid query parameters: " + err.Error(),
		})
		return
	}

	inventories, err := ic.InventoryService.GetInventories(&req)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var responses []models.InventoryResponse
	for _, inv := range inventories {
		responses = append(responses, models.InventoryResponse{
			ID:        inv.ID,
			ProductID: inv.ProductID,
			Product: models.ProductResponse{
				ID:         inv.Product.ID,
				Nama:       inv.Product.Nama,
				Deskripsi:  inv.Product.Deskripsi,
				Harga:      inv.Product.Harga,
				Kategori:   inv.Product.Kategori,
				FotoProduk: inv.Product.FotoProduk,
				CreatedAt:  inv.Product.CreatedAt,
				UpdatedAt:  inv.Product.UpdatedAt,
			},
			Jumlah:    inv.Jumlah,
			Lokasi:    inv.Lokasi,
			CreatedAt: inv.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: inv.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Inventories successfully retrieved",
		Data:    responses,
	})
}

func (ic *InventoryController) GetInventoryByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID is required",
		})
		return
	}

	var idUint uint
	_, err := fmt.Sscanf(idStr, "%d", &idUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID must be a valid number",
		})
		return
	}

	inventory, err := ic.InventoryService.GetInventoryByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := models.InventoryResponse{
		ID:        inventory.ID,
		ProductID: inventory.ProductID,
		Product: models.ProductResponse{
			ID:         inventory.Product.ID,
			Nama:       inventory.Product.Nama,
			Deskripsi:  inventory.Product.Deskripsi,
			Harga:      inventory.Product.Harga,
			Kategori:   inventory.Product.Kategori,
			FotoProduk: inventory.Product.FotoProduk,
			CreatedAt:  inventory.Product.CreatedAt,
			UpdatedAt:  inventory.Product.UpdatedAt,
		},
		Jumlah:    inventory.Jumlah,
		Lokasi:    inventory.Lokasi,
		CreatedAt: inventory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: inventory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Inventory successfully found",
		Data:    response,
	})
}

func (ic *InventoryController) UpdateInventory(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID is required",
		})
		return
	}

	var idUint uint
	_, err := fmt.Sscanf(idStr, "%d", &idUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID must be a valid number",
		})
		return
	}

	var req models.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid data: " + err.Error(),
		})
		return
	}

	inventory, err := ic.InventoryService.GetInventoryByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	inventory.Jumlah = req.Jumlah
	if req.Lokasi != "" {
		inventory.Lokasi = req.Lokasi
	}

	updatedInventory, err := ic.InventoryService.UpdateInventory(inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := models.InventoryResponse{
		ID:        updatedInventory.ID,
		ProductID: updatedInventory.ProductID,
		Product: models.ProductResponse{
			ID:         updatedInventory.Product.ID,
			Nama:       updatedInventory.Product.Nama,
			Deskripsi:  updatedInventory.Product.Deskripsi,
			Harga:      updatedInventory.Product.Harga,
			Kategori:   updatedInventory.Product.Kategori,
			FotoProduk: updatedInventory.Product.FotoProduk,
			CreatedAt:  updatedInventory.Product.CreatedAt,
			UpdatedAt:  updatedInventory.Product.UpdatedAt,
		},
		Jumlah:    updatedInventory.Jumlah,
		Lokasi:    updatedInventory.Lokasi,
		CreatedAt: updatedInventory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: updatedInventory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Inventory successfully updated",
		Data:    response,
	})
}

func (ic *InventoryController) UpdateStock(c *gin.Context) {
	var req models.UpdateStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid data: " + err.Error(),
		})
		return
	}

	inventory, err := ic.InventoryService.UpdateStock(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := models.InventoryResponse{
		ID:        inventory.ID,
		ProductID: inventory.ProductID,
		Product: models.ProductResponse{
			ID:         inventory.Product.ID,
			Nama:       inventory.Product.Nama,
			Deskripsi:  inventory.Product.Deskripsi,
			Harga:      inventory.Product.Harga,
			Kategori:   inventory.Product.Kategori,
			FotoProduk: inventory.Product.FotoProduk,
			CreatedAt:  inventory.Product.CreatedAt,
			UpdatedAt:  inventory.Product.UpdatedAt,
		},
		Jumlah:    inventory.Jumlah,
		Lokasi:    inventory.Lokasi,
		CreatedAt: inventory.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: inventory.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Stock successfully updated",
		Data:    response,
	})
}

func (ic *InventoryController) DeleteInventory(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID is required",
		})
		return
	}

	var idUint uint
	_, err := fmt.Sscanf(idStr, "%d", &idUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "ID must be a valid number",
		})
		return
	}

	err = ic.InventoryService.DeleteInventory(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Inventory successfully deleted",
	})
}
