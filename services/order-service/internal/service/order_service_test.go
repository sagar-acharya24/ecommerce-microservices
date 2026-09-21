package service

import (
	"errors"
	"testing"

	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/client"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/repository"
)

type mockOrderRepository struct {
	orders []model.Order
}

func (m *mockOrderRepository) Create(order *model.Order) error {
	order.ID = uint(len(m.orders) + 1)
	m.orders = append(m.orders, *order)
	return nil
}

func (m *mockOrderRepository) GetByID(id uint) (*model.Order, error) {
	for i := range m.orders {
		if m.orders[i].ID == id {
			order := m.orders[i]
			return &order, nil
		}
	}

	return nil, errors.New("order not found")
}

func (m *mockOrderRepository) GetByUserID(
	userID uint,
) ([]model.Order, error) {
	var result []model.Order

	for _, order := range m.orders {
		if order.UserID == userID {
			result = append(result, order)
		}
	}

	return result, nil
}

func (m *mockOrderRepository) Update(order *model.Order) error {
	for i := range m.orders {
		if m.orders[i].ID == order.ID {
			m.orders[i].Status = order.Status
			return nil
		}
	}

	return errors.New("order not found")
}

var _ repository.OrderRepository = (*mockOrderRepository)(nil)

type mockProductClient struct {
	product *client.Product
	err     error
}

func (m *mockProductClient) GetProduct(
	productID uint,
) (*client.Product, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.product, nil
}

var _ client.ProductClient = (*mockProductClient)(nil)

func TestOrderServiceCreateOrder(t *testing.T) {
	repo := &mockOrderRepository{}

	productClient := &mockProductClient{
		product: &client.Product{
			ID:    35,
			Name:  "Test Laptop",
			Price: 50000,
			Stock: 10,
		},
	}

	orderService := NewOrderService(repo, productClient)

	req := dto.CreateOrderRequest{
		UserID:    1,
		ProductID: 35,
		Quantity:  2,
	}

	order, err := orderService.CreateOrder(req)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	if order.ID == 0 {
		t.Fatal("expected order ID to be generated")
	}

	if order.TotalPrice != 100000 {
		t.Fatalf(
			"expected total price 100000, got %.2f",
			order.TotalPrice,
		)
	}

	if order.Status != model.OrderStatusPending {
		t.Fatalf(
			"expected status %s, got %s",
			model.OrderStatusPending,
			order.Status,
		)
	}
}

func TestOrderServiceCreateOrderInsufficientStock(t *testing.T) {
	repo := &mockOrderRepository{}

	productClient := &mockProductClient{
		product: &client.Product{
			ID:    35,
			Name:  "Test Laptop",
			Price: 50000,
			Stock: 2,
		},
	}

	orderService := NewOrderService(repo, productClient)

	req := dto.CreateOrderRequest{
		UserID:    1,
		ProductID: 35,
		Quantity:  5,
	}

	_, err := orderService.CreateOrder(req)

	if err == nil {
		t.Fatal("expected insufficient stock error")
	}
}

func TestOrderServiceCreateOrderProductError(t *testing.T) {
	repo := &mockOrderRepository{}

	productClient := &mockProductClient{
		err: errors.New("product service unavailable"),
	}

	orderService := NewOrderService(repo, productClient)

	req := dto.CreateOrderRequest{
		UserID:    1,
		ProductID: 35,
		Quantity:  2,
	}

	_, err := orderService.CreateOrder(req)

	if err == nil {
		t.Fatal("expected product service error")
	}
}

func TestOrderServiceCreateOrderInvalidQuantity(t *testing.T) {
	repo := &mockOrderRepository{}

	productClient := &mockProductClient{
		product: &client.Product{
			ID:    35,
			Name:  "Test Laptop",
			Price: 50000,
			Stock: 10,
		},
	}

	orderService := NewOrderService(repo, productClient)

	req := dto.CreateOrderRequest{
		UserID:    1,
		ProductID: 35,
		Quantity:  0,
	}

	_, err := orderService.CreateOrder(req)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestOrderServiceGetOrderByID(t *testing.T) {
	repo := &mockOrderRepository{
		orders: []model.Order{
			{
				ID:         1,
				UserID:     1,
				ProductID:  35,
				Quantity:   2,
				TotalPrice: 100000,
				Status:     model.OrderStatusPending,
			},
		},
	}

	productClient := &mockProductClient{}

	orderService := NewOrderService(repo, productClient)

	order, err := orderService.GetOrderByID(1)
	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}

	if order.ID != 1 {
		t.Fatalf("expected order ID 1, got %d", order.ID)
	}
}

func TestOrderServiceCancelOrder(t *testing.T) {
	repo := &mockOrderRepository{
		orders: []model.Order{
			{
				ID:         1,
				UserID:     1,
				ProductID:  35,
				Quantity:   2,
				TotalPrice: 100000,
				Status:     model.OrderStatusPending,
			},
		},
	}

	productClient := &mockProductClient{}

	orderService := NewOrderService(repo, productClient)

	order, err := orderService.CancelOrder(1)
	if err != nil {
		t.Fatalf("failed to cancel order: %v", err)
	}

	if order.Status != model.OrderStatusCancelled {
		t.Fatalf(
			"expected status %s, got %s",
			model.OrderStatusCancelled,
			order.Status,
		)
	}
}
