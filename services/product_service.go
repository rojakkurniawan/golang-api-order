package services

import (
	"errors"
	"fmt"
	"golang-api/models"
	"io"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (ps *ProductService) CreateProduct(product *models.Product, src io.Reader) (*models.Product, error) {
	const uploadDir = "uploads"

	if err := ps.DB.Create(product).Error; err != nil {
		return nil, errors.New("error creating product " + err.Error())
	}

	if err := os.MkdirAll(uploadDir, 0775); err != nil {
		return nil, fmt.Errorf("error creating directory: %v", err)
	}

	if product.FotoProduk == "" || src == nil {
		if err := ps.DB.Create(product).Error; err != nil {
			return nil, errors.New("error creating product " + err.Error())
		}
		return product, nil
	}

	cleanFilename := filepath.Base(product.FotoProduk)
	pathFolder := filepath.Join(uploadDir, cleanFilename)

	dst, err := os.Create(pathFolder)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("error saving file: %v", err)
	}

	return product, nil
}

func (ps *ProductService) GetProduct(productQuery *models.GetProductRequest) ([]models.Product, error) {
	var product []models.Product

	query := ps.DB

	if productQuery.Kategori != "" {
		query = query.Where("kategori = ?", productQuery.Kategori)
	}

	err := query.Limit(productQuery.Limit).Offset(productQuery.Offset).Find(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	if len(product) == 0 {
		return nil, errors.New("no products found")
	}

	return product, nil
}

func (ps *ProductService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product

	err := ps.DB.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (ps *ProductService) UpdateProduct(product *models.Product, src io.Reader) (*models.Product, error) {
	const uploadDir = "uploads"

	if err := os.MkdirAll(uploadDir, 0775); err != nil {
		return nil, fmt.Errorf("error creating directory: %v", err)
	}

	if product.FotoProduk == "" || src == nil {
		if err := ps.DB.Save(product).Error; err != nil {
			return nil, errors.New("error updating product " + err.Error())
		}
		return product, nil
	}

	cleanFilename := filepath.Base(product.FotoProduk)
	pathFolder := filepath.Join(uploadDir, cleanFilename)

	dst, err := os.Create(pathFolder)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("error saving file: %v", err)
	}

	if err := ps.DB.Save(product).Error; err != nil {
		return nil, errors.New("error updating product " + err.Error())
	}

	return product, nil
}

func (ps *ProductService) DeleteProduct(id uint) error {
	var product models.Product

	err := ps.DB.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	if err := ps.DB.Delete(&product).Error; err != nil {
		return errors.New("error deleting product " + err.Error())
	}

	if product.FotoProduk != "" {
		filePath := filepath.Join("uploads", product.FotoProduk)
		if _, err := os.Stat(filePath); err == nil {
			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("error deleting file: %v", err)
			}
		}
	}

	return nil
}
