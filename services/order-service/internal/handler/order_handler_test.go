package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/client"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/middleware"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/repository"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/service"
	"gorm.io/gorm"
)

type mockHandlerOrderRepository struct {
	orders []model.Order
}

func (m *mockHandlerOrderRepository) Create(order *model.Order) error {
	order.ID = uint(len(m.orders) + 1)
	m.orders = append(m.orders, *order)
	return nil
}

func (m *mockHandlerOrderRepository) GetByID(id uint) (*model.Order, error) {
	for i := range m.orders {
		if m.orders[i].ID == id {
			order := m.orders[i]
			return &order, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (m *mockHandlerOrderRepository) GetByUserID(
	userID uint,
) ([]model.Order, error) {
	var orders []model.Order

	for _, order := range m.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (m *mockHandlerOrderRepository) Update(order *model.Order) error {
	for i := range m.orders {
		if m.orders[i].ID == order.ID {
			m.orders[i].Status = order.Status
			return nil
		}
	}

	return gorm.ErrRecordNotFound
}

var _ repository.OrderRepository = (*mockHandlerOrderRepository)(nil)

type mockHandlerProductClient struct {
	product *client.Product
	err     error
}

func (m *mockHandlerProductClient) GetProduct(
	productID uint,
) (*client.Product, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.product, nil
}

var _ client.ProductClient = (*mockHandlerProductClient)(nil)

func setupTestRouter(
	orderService service.OrderService,
) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	handler := NewOrderHandler(orderService)

	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, uint(1))
		c.Next()
	})

	router.POST("/api/v1/orders", handler.CreateOrder)
	router.GET("/api/v1/orders/:id", handler.GetOrder)
	router.GET(
		"/api/v1/user/:user_id/orders",
		handler.GetUserOrders,
	)
	router.PUT(
		"/api/v1/orders/:id/cancel",
		handler.CancelOrder,
	)

	return router
}

func TestCreateOrderHandler(t *testing.T) {
	repo := &mockHandlerOrderRepository{}

	productClient := &mockHandlerProductClient{
		product: &client.Product{
			ID:    35,
			Name:  "Test Laptop",
			Price: 50000,
			Stock: 10,
		},
	}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	body := `{
		"user_id": 999,
		"product_id": 35,
		"quantity": 2
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/orders",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	if len(repo.orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(repo.orders))
	}

	if repo.orders[0].UserID != 1 {
		t.Fatalf(
			"expected order user_id 1, got %d",
			repo.orders[0].UserID,
		)
	}
}

func TestGetOrderHandlerForbidden(t *testing.T) {
	repo := &mockHandlerOrderRepository{
		orders: []model.Order{
			{
				ID:         1,
				UserID:     2, // Different from authenticated user (1)
				ProductID:  35,
				Quantity:   2,
				TotalPrice: 100000,
				Status:     model.OrderStatusPending,
			},
		},
	}

	productClient := &mockHandlerProductClient{}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/orders/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			recorder.Code,
		)
	}
}

func TestCreateOrderHandlerInvalidRequest(t *testing.T) {
	repo := &mockHandlerOrderRepository{}

	productClient := &mockHandlerProductClient{
		product: &client.Product{
			ID:    35,
			Name:  "Test Laptop",
			Price: 50000,
			Stock: 10,
		},
	}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	body := `{
		"user_id": 1,
		"product_id": 35,
		"quantity": 0
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/orders",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestGetOrderHandler(t *testing.T) {
	repo := &mockHandlerOrderRepository{
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

	productClient := &mockHandlerProductClient{}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/orders/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestGetOrderHandlerNotFound(t *testing.T) {
	repo := &mockHandlerOrderRepository{}

	productClient := &mockHandlerProductClient{}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/orders/999",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}
}

func TestGetUserOrdersHandler(t *testing.T) {
	repo := &mockHandlerOrderRepository{
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

	productClient := &mockHandlerProductClient{}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/user/1/orders",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestCancelOrderHandler(t *testing.T) {
	repo := &mockHandlerOrderRepository{
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

	productClient := &mockHandlerProductClient{}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/orders/1/cancel",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestCancelOrderHandlerForbidden(t *testing.T) {
	repo := &mockHandlerOrderRepository{
		orders: []model.Order{
			{
				ID:         1,
				UserID:     2, // Order belongs to another user.
				ProductID:  35,
				Quantity:   2,
				TotalPrice: 100000,
				Status:     model.OrderStatusPending,
			},
		},
	}

	productClient := &mockHandlerProductClient{}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/orders/1/cancel",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			recorder.Code,
		)
	}

	// Make sure the order was NOT cancelled.
	if repo.orders[0].Status != model.OrderStatusPending {
		t.Fatalf(
			"expected order status to remain %s, got %s",
			model.OrderStatusPending,
			repo.orders[0].Status,
		)
	}
}

func TestCancelOrderHandlerAlreadyCancelled(t *testing.T) {
	repo := &mockHandlerOrderRepository{
		orders: []model.Order{
			{
				ID:         1,
				UserID:     1,
				ProductID:  35,
				Quantity:   2,
				TotalPrice: 100000,
				Status:     model.OrderStatusCancelled,
			},
		},
	}

	productClient := &mockHandlerProductClient{}

	orderService := service.NewOrderService(
		repo,
		productClient,
	)

	router := setupTestRouter(orderService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/orders/1/cancel",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}
