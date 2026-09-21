package service

import (
	"fmt"

	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/client"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(req dto.CreateOrderRequest) (*model.Order, error)
	GetOrderByID(id uint) (*model.Order, error)
	GetOrdersByUserID(userID uint) ([]model.Order, error)
	CancelOrder(id uint) (*model.Order, error)
}

type orderService struct {
	repo          repository.OrderRepository
	productClient client.ProductClient
}

func NewOrderService(
	repo repository.OrderRepository,
	productClient client.ProductClient,
) OrderService {
	return &orderService{
		repo:          repo,
		productClient: productClient,
	}
}

func (s *orderService) CreateOrder(
	req dto.CreateOrderRequest,
) (*model.Order, error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user_id must be greater than 0")
	}

	if req.ProductID == 0 {
		return nil, fmt.Errorf("product_id must be greater than 0")
	}

	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	product, err := s.productClient.GetProduct(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get product: %w",
			err,
		)
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf(
			"insufficient stock: available=%d requested=%d",
			product.Stock,
			req.Quantity,
		)
	}

	totalPrice := product.Price * float64(req.Quantity)

	order := &model.Order{
		UserID:     req.UserID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPending,
	}

	if err := s.repo.Create(order); err != nil {
		return nil, fmt.Errorf(
			"failed to create order: %w",
			err,
		)
	}

	return order, nil
}

func (s *orderService) GetOrderByID(id uint) (*model.Order, error) {
	if id == 0 {
		return nil, fmt.Errorf("order id must be greater than 0")
	}

	order, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrdersByUserID(
	userID uint,
) ([]model.Order, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user id must be greater than 0")
	}

	orders, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *orderService) CancelOrder(id uint) (*model.Order, error) {
	if id == 0 {
		return nil, fmt.Errorf("order id must be greater than 0")
	}

	order, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if order.Status == model.OrderStatusCancelled {
		return nil, fmt.Errorf("order is already cancelled")
	}

	order.Status = model.OrderStatusCancelled

	if err := s.repo.Update(order); err != nil {
		return nil, fmt.Errorf(
			"failed to cancel order: %w",
			err,
		)
	}

	updatedOrder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}
