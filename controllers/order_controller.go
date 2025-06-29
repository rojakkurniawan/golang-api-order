package controllers

import (
	"fmt"
	"golang-api/models"
	"golang-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	OrderService *services.OrderService
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		OrderService: services.NewOrderService(db),
	}
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid data: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "user not authenticated",
		})
		return
	}

	order, err := oc.OrderService.CreateOrder(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := oc.convertToOrderResponse(order)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Order successfully created",
		Data:    response,
	})
}

func (oc *OrderController) GetOrders(c *gin.Context) {
	var req models.GetOrderRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid query parameters: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "user not authenticated",
		})
		return
	}

	req.UserID = userID.(uint)

	orders, err := oc.OrderService.GetOrders(&req)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var responses []models.OrderResponse
	for _, order := range orders {
		responses = append(responses, oc.convertToOrderResponse(&order))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Orders successfully retrieved",
		Data:    responses,
	})
}

func (oc *OrderController) GetOrderByID(c *gin.Context) {
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

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "user not authenticated",
		})
		return
	}

	order, err := oc.OrderService.GetOrderByIDAndUserID(idUint, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := oc.convertToOrderResponse(order)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order successfully found",
		Data:    response,
	})
}

func (oc *OrderController) UpdateOrderStatus(c *gin.Context) {
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

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid data: " + err.Error(),
		})
		return
	}

	order, err := oc.OrderService.UpdateOrderStatus(idUint, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := oc.convertToOrderResponse(order)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order status successfully updated",
		Data:    response,
	})
}

func (oc *OrderController) DeleteOrder(c *gin.Context) {
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

	err = oc.OrderService.DeleteOrder(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order successfully deleted",
	})
}

func (oc *OrderController) convertToOrderResponse(order *models.Order) models.OrderResponse {
	var orderItems []models.OrderItemResponse
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, models.OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Product: models.ProductResponse{
				ID:         item.Product.ID,
				Nama:       item.Product.Nama,
				Deskripsi:  item.Product.Deskripsi,
				Harga:      item.Product.Harga,
				Kategori:   item.Product.Kategori,
				FotoProduk: item.Product.FotoProduk,
				CreatedAt:  item.Product.CreatedAt,
				UpdatedAt:  item.Product.UpdatedAt,
			},
			Jumlah:   item.Jumlah,
			Harga:    item.Harga,
			Subtotal: item.Subtotal,
		})
	}

	return models.OrderResponse{
		ID:     order.ID,
		UserID: order.UserID,
		User: models.UserResponse{
			ID:    order.User.ID,
			Name:  order.User.Name,
			Email: order.User.Email,
		},
		TotalHarga:   order.TotalHarga,
		Status:       order.Status,
		TanggalOrder: order.TanggalOrder.Format("2006-01-02 15:04:05"),
		OrderItems:   orderItems,
		CreatedAt:    order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
