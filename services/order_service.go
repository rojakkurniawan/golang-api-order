package services

import (
	"errors"
	"golang-api/models"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{DB: db}
}

func (os *OrderService) CreateOrder(userID uint, req *models.CreateOrderRequest) (*models.Order, error) {
	tx := os.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	order := models.Order{
		UserID:       userID,
		Status:       "pending",
		TanggalOrder: time.Now(),
		TotalHarga:   0,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("error creating order: " + err.Error())
	}

	var totalHarga float64
	var orderItems []models.OrderItem

	for _, item := range req.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("product not found")
			}
			return nil, err
		}

		subtotal := product.Harga * float64(item.Jumlah)

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Jumlah:    item.Jumlah,
			Harga:     product.Harga,
			Subtotal:  subtotal,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("error creating order item: " + err.Error())
		}

		orderItems = append(orderItems, orderItem)
		totalHarga += subtotal
	}

	order.TotalHarga = totalHarga
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("error updating order total: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("error committing transaction: " + err.Error())
	}

	if err := os.DB.Preload("User").Preload("OrderItems.Product").First(&order, order.ID).Error; err != nil {
		return nil, errors.New("error loading order with relations")
	}

	return &order, nil
}

func (os *OrderService) GetOrders(req *models.GetOrderRequest) ([]models.Order, error) {
	var orders []models.Order

	query := os.DB.Preload("User").Preload("OrderItems.Product")

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	err := query.Order("created_at DESC").Limit(req.Limit).Offset(req.Offset).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}

	return orders, nil
}

func (os *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order

	err := os.DB.Preload("User").Preload("OrderItems.Product").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (os *OrderService) GetOrderByIDAndUserID(id uint, userID uint) (*models.Order, error) {
	var order models.Order

	err := os.DB.Preload("User").Preload("OrderItems.Product").Where("id = ? AND user_id = ?", id, userID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (os *OrderService) UpdateOrderStatus(id uint, status string) (*models.Order, error) {
	var order models.Order

	if err := os.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	order.Status = status
	if err := os.DB.Save(&order).Error; err != nil {
		return nil, errors.New("error updating order status: " + err.Error())
	}

	if err := os.DB.Preload("User").Preload("OrderItems.Product").First(&order, order.ID).Error; err != nil {
		return nil, errors.New("error loading order with relations")
	}

	return &order, nil
}

func (os *OrderService) DeleteOrder(id uint) error {
	tx := os.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var order models.Order
	if err := tx.First(&order, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return err
	}

	if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return errors.New("error deleting order items: " + err.Error())
	}

	if err := tx.Delete(&order).Error; err != nil {
		tx.Rollback()
		return errors.New("error deleting order: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("error committing transaction: " + err.Error())
	}

	return nil
}
