package controllers

import (
	"fmt"
	"golang-api/models"
	"golang-api/services"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductController struct {
	ProductService *services.ProductService
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{
		ProductService: services.NewProductService(db),
	}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req models.AddProductRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid data",
		})
		return
	}

	file, err := c.FormFile("foto_produk")
	var fileName string
	var src multipart.File

	if err == nil {
		src, err = file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "error opening file",
			})
			return
		}
		defer src.Close()

		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "error reading file",
			})
			return
		}
		filetype := http.DetectContentType(buffer)
		if !(filetype == "image/jpeg" || filetype == "image/png" || filetype == "image/gif") {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "file must be image jpg, png, gif",
			})
			return
		}
		src.Seek(0, 0)
		fileName = fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(file.Filename))

	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "error file not valid",
		})
		return
	}

	newProduct := models.Product{
		Nama:       req.Nama,
		Deskripsi:  req.Deskripsi,
		Harga:      float64(req.Harga),
		Kategori:   req.Kategori,
		FotoProduk: fileName,
	}

	product, err := pc.ProductService.CreateProduct(&newProduct, src)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if product.FotoProduk != "" {
		fotoLink := fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/"+product.FotoProduk)
		product.FotoProduk = fotoLink
	} else {
		product.FotoProduk = fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/default.png")
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product successfully created",
		Data: models.ProductResponse{
			ID:         product.ID,
			Nama:       product.Nama,
			Deskripsi:  product.Deskripsi,
			Harga:      product.Harga,
			Kategori:   product.Kategori,
			FotoProduk: product.FotoProduk,
			CreatedAt:  product.CreatedAt,
			UpdatedAt:  product.CreatedAt,
		},
	})

}

func (pc *ProductController) GetProduct(c *gin.Context) {
	var req models.GetProductRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid query paramater",
			Data:    err.Error(),
		})
		return
	}

	products, err := pc.ProductService.GetProduct(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var responses []models.ProductResponse

	for _, p := range products {
		if p.FotoProduk != "" {
			fotoLink := fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/"+p.FotoProduk)
			p.FotoProduk = fotoLink
		} else {
			p.FotoProduk = fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/default.png")
		}

		responses = append(responses, models.ProductResponse{
			ID:         p.ID,
			Nama:       p.Nama,
			Deskripsi:  p.Deskripsi,
			Harga:      p.Harga,
			Kategori:   p.Kategori,
			FotoProduk: p.FotoProduk,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product successfully retrieved",
		Data:    responses,
	})
}

func (pc *ProductController) GetProductByID(c *gin.Context) {
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

	product, err := pc.ProductService.GetProductByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if product.FotoProduk != "" {
		fotoLink := fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/"+product.FotoProduk)
		product.FotoProduk = fotoLink
	} else {
		product.FotoProduk = fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/default.png")
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product successfully found",
		Data: models.ProductResponse{
			ID:         product.ID,
			Nama:       product.Nama,
			Deskripsi:  product.Deskripsi,
			Harga:      product.Harga,
			Kategori:   product.Kategori,
			FotoProduk: product.FotoProduk,
			CreatedAt:  product.CreatedAt,
			UpdatedAt:  product.UpdatedAt,
		},
	})
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var req models.UpdateProductRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid data",
		})
		return
	}

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

	file, err := c.FormFile("foto_produk")
	var fileName string
	var src multipart.File

	if err == nil {
		src, err = file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "error opening file",
			})
			return
		}
		defer src.Close()

		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "error reading file",
			})
			return
		}
		filetype := http.DetectContentType(buffer)
		if !(filetype == "image/jpeg" || filetype == "image/png" || filetype == "image/gif") {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "file must be image jpg, png, gif",
			})
			return
		}
		src.Seek(0, 0)
		fileName = fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(file.Filename))

	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "error file not valid",
		})
		return
	}

	product, err := pc.ProductService.GetProductByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if product.FotoProduk != "" && fileName != "" {
		oldFilePath := filepath.Join("uploads", product.FotoProduk)
		if _, err := os.Stat(oldFilePath); err == nil {
			err := os.Remove(oldFilePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "error deleting old file",
				})
				return
			}
		}
	}

	product.Nama = req.Nama
	product.Harga = float64(req.Harga)
	product.Kategori = req.Kategori
	product.Deskripsi = req.Deskripsi
	if fileName != "" {
		product.FotoProduk = fileName
	}

	updatedProduct, err := pc.ProductService.UpdateProduct(product, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	if updatedProduct.FotoProduk != "" {
		fotoLink := fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/"+updatedProduct.FotoProduk)
		updatedProduct.FotoProduk = fotoLink
	} else {
		updatedProduct.FotoProduk = fmt.Sprintf("http://%s/%s", c.Request.Host, "api/products/images/default.png")
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product successfully updated",
		Data: models.ProductResponse{
			ID:         updatedProduct.ID,
			Nama:       updatedProduct.Nama,
			Deskripsi:  updatedProduct.Deskripsi,
			Harga:      updatedProduct.Harga,
			Kategori:   updatedProduct.Kategori,
			FotoProduk: updatedProduct.FotoProduk,
			CreatedAt:  updatedProduct.CreatedAt,
			UpdatedAt:  updatedProduct.UpdatedAt,
		},
	})
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
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

	err = pc.ProductService.DeleteProduct(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product successfully deleted",
	})
}
