package repository

import (
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/config"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/database"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/model"
)

func setupTestRepository(t *testing.T) (ProductRepository, func()) {
	t.Helper()

	// The test runs from the repository package directory,
	// so explicitly load the root .env file.
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

	repo := NewProductRepository(db)

	cleanup := func() {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
			Delete(&model.Product{}).Error; err != nil {
			t.Errorf("failed to clean up test products: %v", err)
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

func TestProductRepositoryCreate(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	product := &model.Product{
		Name:        "Test MacBook",
		Description: "Test laptop",
		Price:       99999,
		Stock:       5,
	}

	err := repo.Create(product)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if product.ID == 0 {
		t.Fatal("expected product ID to be generated")
	}

	t.Logf("Created product with ID: %d", product.ID)
}

func TestProductRepositoryGetByID(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	product := &model.Product{
		Name:        "Test Keyboard",
		Description: "Mechanical keyboard",
		Price:       4999,
		Stock:       20,
	}

	if err := repo.Create(product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	foundProduct, err := repo.GetByID(product.ID)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if foundProduct.ID != product.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			product.ID,
			foundProduct.ID,
		)
	}

	if foundProduct.Name != "Test Keyboard" {
		t.Fatalf(
			"expected name %q, got %q",
			"Test Keyboard",
			foundProduct.Name,
		)
	}

	t.Logf(
		"Found product: ID=%d Name=%s Price=%.2f Stock=%d",
		foundProduct.ID,
		foundProduct.Name,
		foundProduct.Price,
		foundProduct.Stock,
	)
}

func TestProductRepositoryGetByIDNotFound(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	_, err := repo.GetByID(999999)

	if err == nil {
		t.Fatal("expected error for non-existing product")
	}

	if err != gorm.ErrRecordNotFound {
		t.Fatalf(
			"expected gorm.ErrRecordNotFound, got %v",
			err,
		)
	}

	t.Log("Correctly returned record not found")
}

func TestProductRepositoryGetAll(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	product1 := &model.Product{
		Name:  "Test Mouse",
		Price: 1999,
		Stock: 15,
	}

	product2 := &model.Product{
		Name:  "Test Monitor",
		Price: 24999,
		Stock: 8,
	}

	if err := repo.Create(product1); err != nil {
		t.Fatalf("failed to create product1: %v", err)
	}

	if err := repo.Create(product2); err != nil {
		t.Fatalf("failed to create product2: %v", err)
	}

	products, err := repo.GetAll()
	if err != nil {
		t.Fatalf("failed to get products: %v", err)
	}

	if len(products) < 2 {
		t.Fatalf(
			"expected at least 2 products, got %d",
			len(products),
		)
	}

	t.Logf("Found %d products", len(products))
}

func TestProductRepositoryUpdate(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	product := &model.Product{
		Name:        "Test Headphones",
		Description: "Original description",
		Price:       3999,
		Stock:       10,
	}

	if err := repo.Create(product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	product.Name = "Updated Headphones"
	product.Price = 4499
	product.Stock = 20

	if err := repo.Update(product); err != nil {
		t.Fatalf("failed to update product: %v", err)
	}

	updatedProduct, err := repo.GetByID(product.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated product: %v", err)
	}

	if updatedProduct.Name != "Updated Headphones" {
		t.Fatalf(
			"expected name %q, got %q",
			"Updated Headphones",
			updatedProduct.Name,
		)
	}

	if updatedProduct.Price != 4499 {
		t.Fatalf(
			"expected price %.2f, got %.2f",
			4499.0,
			updatedProduct.Price,
		)
	}

	if updatedProduct.Stock != 20 {
		t.Fatalf(
			"expected stock %d, got %d",
			20,
			updatedProduct.Stock,
		)
	}

	t.Logf(
		"Updated product: ID=%d Name=%s Price=%.2f Stock=%d",
		updatedProduct.ID,
		updatedProduct.Name,
		updatedProduct.Price,
		updatedProduct.Stock,
	)
}

func TestProductRepositoryDelete(t *testing.T) {
	repo, cleanup := setupTestRepository(t)
	defer cleanup()

	product := &model.Product{
		Name:  "Test Delete Product",
		Price: 999,
		Stock: 3,
	}

	if err := repo.Create(product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	productID := product.ID

	if err := repo.Delete(productID); err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}

	_, err := repo.GetByID(productID)

	if err != gorm.ErrRecordNotFound {
		t.Fatalf(
			"expected gorm.ErrRecordNotFound after delete, got %v",
			err,
		)
	}

	t.Logf("Successfully deleted product with ID: %d", productID)
}
