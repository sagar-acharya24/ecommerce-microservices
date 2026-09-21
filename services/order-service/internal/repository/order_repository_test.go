package repository

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/config"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/database"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/model"
	"gorm.io/gorm"
)

func setupTestRepository(t *testing.T) (OrderRepository, func()) {
	t.Helper()

	if err := godotenv.Load("../../../../.env"); err != nil {
		t.Fatalf("failed to load .env: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	repo := NewOrderRepository(db)

	cleanup := func() {
		if err := db.Session(&gorm.Session{
			AllowGlobalUpdate: true,
		}).Delete(&model.Order{}).Error; err != nil {
			t.Errorf("failed to clean up test orders: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			t.Errorf("failed to get underlying database connection: %v", err)
			return
		}

		if err := sqlDB.Close(); err != nil {
			t.Errorf("failed to close database connection: %v", err)
		}
	}

	return repo, cleanup
}

func TestOrderRepositoryCreate(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	order := &model.Order{
		UserID:     1,
		ProductID:  35,
		Quantity:   2,
		TotalPrice: 100000,
		Status:     model.OrderStatusPending,
	}

	err := repo.Create(order)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	if order.ID == 0 {
		t.Fatal("expected order ID to be generated")
	}

	if order.Status != model.OrderStatusPending {
		t.Fatalf(
			"expected status %s, got %s",
			model.OrderStatusPending,
			order.Status,
		)
	}
}

func TestOrderRepositoryGetByID(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	order := &model.Order{
		UserID:     1,
		ProductID:  35,
		Quantity:   2,
		TotalPrice: 100000,
		Status:     model.OrderStatusPending,
	}

	if err := repo.Create(order); err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	found, err := repo.GetByID(order.ID)
	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}

	if found.ID != order.ID {
		t.Fatalf(
			"expected order ID %d, got %d",
			order.ID,
			found.ID,
		)
	}

	if found.UserID != 1 {
		t.Fatalf(
			"expected user ID 1, got %d",
			found.UserID,
		)
	}

	if found.ProductID != 35 {
		t.Fatalf(
			"expected product ID 35, got %d",
			found.ProductID,
		)
	}
}

func TestOrderRepositoryGetByUserID(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	orders := []model.Order{
		{
			UserID:     1,
			ProductID:  35,
			Quantity:   2,
			TotalPrice: 100000,
			Status:     model.OrderStatusPending,
		},
		{
			UserID:     1,
			ProductID:  35,
			Quantity:   1,
			TotalPrice: 50000,
			Status:     model.OrderStatusConfirmed,
		},
		{
			UserID:     2,
			ProductID:  35,
			Quantity:   3,
			TotalPrice: 150000,
			Status:     model.OrderStatusPending,
		},
	}

	for i := range orders {
		if err := repo.Create(&orders[i]); err != nil {
			t.Fatalf("failed to create order: %v", err)
		}
	}

	found, err := repo.GetByUserID(1)
	if err != nil {
		t.Fatalf("failed to get user orders: %v", err)
	}

	if len(found) != 2 {
		t.Fatalf(
			"expected 2 orders for user 1, got %d",
			len(found),
		)
	}

	for _, order := range found {
		if order.UserID != 1 {
			t.Fatalf(
				"expected all orders to belong to user 1, got user %d",
				order.UserID,
			)
		}
	}
}

func TestOrderRepositoryUpdate(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	order := &model.Order{
		UserID:     1,
		ProductID:  35,
		Quantity:   2,
		TotalPrice: 100000,
		Status:     model.OrderStatusPending,
	}

	if err := repo.Create(order); err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	order.Status = model.OrderStatusCancelled

	if err := repo.Update(order); err != nil {
		t.Fatalf("failed to update order: %v", err)
	}

	found, err := repo.GetByID(order.ID)
	if err != nil {
		t.Fatalf("failed to get updated order: %v", err)
	}

	if found.Status != model.OrderStatusCancelled {
		t.Fatalf(
			"expected status %s, got %s",
			model.OrderStatusCancelled,
			found.Status,
		)
	}
}
