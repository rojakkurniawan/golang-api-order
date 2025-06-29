package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID       uint        `json:"user_id" gorm:"not null"`
	User         User        `json:"user" gorm:"foreignKey:UserID"`
	TotalHarga   float64     `json:"total_harga" gorm:"default:0"`
	Status       string      `json:"status" gorm:"default:pending"`
	TanggalOrder time.Time   `json:"tanggal_order"`
	OrderItems   []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id" gorm:"not null"`
	Order     Order   `json:"order" gorm:"foreignKey:OrderID"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Jumlah    int     `json:"jumlah" gorm:"not null"`
	Harga     float64 `json:"harga" gorm:"not null"`
	Subtotal  float64 `json:"subtotal" gorm:"not null"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
}

type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Jumlah    int  `json:"jumlah" binding:"required,min=1"`
}

type OrderResponse struct {
	ID           uint                `json:"id"`
	UserID       uint                `json:"user_id"`
	User         UserResponse        `json:"user"`
	TotalHarga   float64             `json:"total_harga"`
	Status       string              `json:"status"`
	TanggalOrder string              `json:"tanggal_order"`
	OrderItems   []OrderItemResponse `json:"order_items"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        uint            `json:"id"`
	OrderID   uint            `json:"order_id"`
	ProductID uint            `json:"product_id"`
	Product   ProductResponse `json:"product"`
	Jumlah    int             `json:"jumlah"`
	Harga     float64         `json:"harga"`
	Subtotal  float64         `json:"subtotal"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetOrderRequest struct {
	Status string `form:"status"`
	UserID uint   `form:"user_id"`
	Limit  int    `form:"limit,default=10"`
	Offset int    `form:"offset,default=0"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}
