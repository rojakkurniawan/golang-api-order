package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Nama       string  `json:"nama"`
	Deskripsi  string  `json:"deskripsi"`
	Harga      float64 `json:"harga"`
	Kategori   string  `json:"kategori"`
	FotoProduk string  `json:"foto_produk"`
}

type AddProductRequest struct {
	Nama      string  `form:"nama" binding:"required"`
	Deskripsi string  `form:"deskripsi"`
	Harga     float64 `form:"harga" binding:"required"`
	Kategori  string  `form:"kategori" binding:"required"`
}

type ProductResponse struct {
	ID         uint      `json:"id"`
	Nama       string    `json:"nama"`
	Deskripsi  string    `json:"deskripsi"`
	Harga      float64   `json:"harga"`
	Kategori   string    `json:"kategori"`
	FotoProduk string    `json:"foto_produk"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetProductRequest struct {
	Kategori string `form:"kategori"`
	Limit    int    `form:"limit,default=10"`
	Offset   int    `form:"offset,default=0"`
}

type UpdateProductRequest struct {
	Nama      string  `form:"nama"`
	Deskripsi string  `form:"deskripsi"`
	Harga     float64 `form:"harga"`
	Kategori  string  `form:"kategori"`
}
