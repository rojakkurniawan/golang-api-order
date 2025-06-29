package services

import (
	"errors"
	"golang-api/models"

	"gorm.io/gorm"
)

type InventoryService struct {
	DB *gorm.DB
}

func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{DB: db}
}

func (is *InventoryService) CreateInventory(inventory *models.Inventory) (*models.Inventory, error) {
	var product models.Product
	if err := is.DB.First(&product, inventory.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	var existingInventory models.Inventory
	err := is.DB.Where("product_id = ? AND lokasi = ?", inventory.ProductID, inventory.Lokasi).First(&existingInventory).Error
	if err == nil {
		return nil, errors.New("inventory for this product and location already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := is.DB.Create(inventory).Error; err != nil {
		return nil, errors.New("error creating inventory: " + err.Error())
	}

	if err := is.DB.Preload("Product").First(inventory, inventory.ID).Error; err != nil {
		return nil, errors.New("error loading inventory with product")
	}

	return inventory, nil
}

func (is *InventoryService) GetInventories(req *models.GetInventoryRequest) ([]models.Inventory, error) {
	var inventories []models.Inventory

	query := is.DB.Preload("Product")

	if req.ProductID != 0 {
		query = query.Where("product_id = ?", req.ProductID)
	}

	if req.Lokasi != "" {
		query = query.Where("lokasi = ?", req.Lokasi)
	}

	err := query.Limit(req.Limit).Offset(req.Offset).Find(&inventories).Error
	if err != nil {
		return nil, err
	}

	if len(inventories) == 0 {
		return nil, errors.New("no inventories found")
	}

	return inventories, nil
}

func (is *InventoryService) GetInventoryByID(id uint) (*models.Inventory, error) {
	var inventory models.Inventory

	err := is.DB.Preload("Product").First(&inventory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory not found")
		}
		return nil, err
	}

	return &inventory, nil
}

func (is *InventoryService) GetInventoryByProductAndLocation(productID uint, lokasi string) (*models.Inventory, error) {
	var inventory models.Inventory

	err := is.DB.Preload("Product").Where("product_id = ? AND lokasi = ?", productID, lokasi).First(&inventory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory not found for this product and location")
		}
		return nil, err
	}

	return &inventory, nil
}

func (is *InventoryService) UpdateInventory(inventory *models.Inventory) (*models.Inventory, error) {
	var product models.Product
	if err := is.DB.First(&product, inventory.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	if err := is.DB.Save(inventory).Error; err != nil {
		return nil, errors.New("error updating inventory: " + err.Error())
	}

	if err := is.DB.Preload("Product").First(inventory, inventory.ID).Error; err != nil {
		return nil, errors.New("error loading inventory with product")
	}

	return inventory, nil
}

func (is *InventoryService) UpdateStock(req *models.UpdateStockRequest) (*models.Inventory, error) {
	inventory, err := is.GetInventoryByProductAndLocation(req.ProductID, req.Lokasi)
	if err != nil {
		return nil, err
	}

	newStock := inventory.Jumlah + req.Jumlah
	if newStock < 0 {
		return nil, errors.New("insufficient stock - operation would result in negative stock")
	}

	inventory.Jumlah = newStock

	if err := is.DB.Save(inventory).Error; err != nil {
		return nil, errors.New("error updating stock: " + err.Error())
	}

	if err := is.DB.Preload("Product").First(inventory, inventory.ID).Error; err != nil {
		return nil, errors.New("error loading inventory with product")
	}

	return inventory, nil
}

func (is *InventoryService) DeleteInventory(id uint) error {
	var inventory models.Inventory

	err := is.DB.First(&inventory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory not found")
		}
		return err
	}

	if err := is.DB.Delete(&inventory).Error; err != nil {
		return errors.New("error deleting inventory: " + err.Error())
	}

	return nil
}
