package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/middleware"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/service"
	"gorm.io/gorm"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authenticated user identity is required",
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authenticated user identity",
		})
		return
	}

	// The authenticated JWT identity is the source of truth.
	// Ignore the user_id supplied by the client.
	req.UserID = userID

	order, err := h.service.CreateOrder(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authenticated user identity is required",
		})
		return
	}

	authenticatedUserID, ok := userIDValue.(uint)
	if !ok || authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authenticated user identity",
		})
		return
	}

	order, err := h.service.GetOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if order.UserID != authenticatedUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not authorized to access this order",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authenticated user identity is required",
		})
		return
	}

	authenticatedUserID, ok := userIDValue.(uint)
	if !ok || authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authenticated user identity",
		})
		return
	}

	requestedUserID, err := strconv.ParseUint(
		c.Param("user_id"),
		10,
		64,
	)

	if err != nil || requestedUserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	if uint(requestedUserID) != authenticatedUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not authorized to access these orders",
		})
		return
	}

	orders, err := h.service.GetOrdersByUserID(uint(requestedUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authenticated user identity is required",
		})
		return
	}

	authenticatedUserID, ok := userIDValue.(uint)
	if !ok || authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authenticated user identity",
		})
		return
	}

	// First fetch the order so we can verify ownership
	// before performing the cancellation.
	existingOrder, err := h.service.GetOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if existingOrder.UserID != authenticatedUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not authorized to cancel this order",
		})
		return
	}

	order, err := h.service.CancelOrder(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}
