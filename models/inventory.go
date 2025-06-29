package models

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Jumlah    int     `json:"jumlah" gorm:"default:0"`
	Lokasi    string  `json:"lokasi" gorm:"not null"`
}

type AddInventoryRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Jumlah    int    `json:"jumlah" binding:"required,min=0"`
	Lokasi    string `json:"lokasi" binding:"required"`
}

type UpdateInventoryRequest struct {
	Jumlah int    `json:"jumlah" binding:"required,min=0"`
	Lokasi string `json:"lokasi"`
}

type UpdateStockRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Jumlah    int    `json:"jumlah" binding:"required"`
	Lokasi    string `json:"lokasi" binding:"required"`
}

type InventoryResponse struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"product_id"`
	Product   ProductResponse `json:"product"`
	Jumlah    int             `json:"jumlah"`
	Lokasi    string          `json:"lokasi"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

type GetInventoryRequest struct {
	ProductID uint   `form:"product_id"`
	Lokasi    string `form:"lokasi"`
	Limit     int    `form:"limit,default=10"`
	Offset    int    `form:"offset,default=0"`
}
