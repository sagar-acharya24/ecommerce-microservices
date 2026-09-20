package service

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/model"
)

type mockProductCache struct {
	products  map[uint]*model.Product
	getErr    error
	setErr    error
	deleteErr error
}

func newMockProductCache() *mockProductCache {
	return &mockProductCache{
		products: make(map[uint]*model.Product),
	}
}

func (m *mockProductCache) GetProduct(
	_ context.Context,
	id uint,
) (*model.Product, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}

	product, exists := m.products[id]
	if !exists {
		return nil, nil
	}

	return product, nil
}

func (m *mockProductCache) SetProduct(
	_ context.Context,
	product *model.Product,
) error {
	if m.setErr != nil {
		return m.setErr
	}

	m.products[product.ID] = product
	return nil
}

func (m *mockProductCache) DeleteProduct(
	_ context.Context,
	id uint,
) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}

	delete(m.products, id)
	return nil
}

type mockProductRepository struct {
	products map[uint]*model.Product
	nextID   uint
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[uint]*model.Product),
		nextID:   1,
	}
}

func (m *mockProductRepository) Create(product *model.Product) error {
	product.ID = m.nextID
	m.nextID++

	m.products[product.ID] = product

	return nil
}

func (m *mockProductRepository) GetByID(id uint) (*model.Product, error) {
	product, exists := m.products[id]

	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	return product, nil
}

func (m *mockProductRepository) GetAll() ([]model.Product, error) {
	products := make([]model.Product, 0, len(m.products))

	for _, product := range m.products {
		products = append(products, *product)
	}

	return products, nil
}

func (m *mockProductRepository) Update(product *model.Product) error {
	if _, exists := m.products[product.ID]; !exists {
		return gorm.ErrRecordNotFound
	}

	m.products[product.ID] = product

	return nil
}

func (m *mockProductRepository) Delete(id uint) error {
	if _, exists := m.products[id]; !exists {
		return gorm.ErrRecordNotFound
	}

	delete(m.products, id)

	return nil
}

func TestProductServiceCreateProduct(t *testing.T) {
	repo := newMockProductRepository()
	cache := newMockProductCache()

	service := NewProductService(
		repo,
		cache,
	)

	product := &model.Product{
		Name:        "MacBook Pro",
		Description: "Apple laptop",
		Price:       149999,
		Stock:       10,
	}

	err := service.CreateProduct(product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID == 0 {
		t.Fatal("expected product ID to be generated")
	}

	t.Logf("Created product with ID: %d", product.ID)
}

func TestProductServiceCreateProductValidation(t *testing.T) {
	repo := newMockProductRepository()
	cache := newMockProductCache()

	service := NewProductService(
		repo,
		cache,
	)

	tests := []struct {
		name    string
		product model.Product
	}{
		{
			name: "missing name",
			product: model.Product{
				Price: 1000,
				Stock: 10,
			},
		},
		{
			name: "invalid price",
			product: model.Product{
				Name:  "Test Product",
				Price: 0,
				Stock: 10,
			},
		},
		{
			name: "negative stock",
			product: model.Product{
				Name:  "Test Product",
				Price: 1000,
				Stock: -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateProduct(&tt.product)

			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestProductServiceGetProductByID(t *testing.T) {
	repo := newMockProductRepository()
	cache := newMockProductCache()

	service := NewProductService(
		repo,
		cache,
	)
	product := &model.Product{
		Name:  "iPhone",
		Price: 79999,
		Stock: 5,
	}

	if err := service.CreateProduct(product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	found, err := service.GetProductByID(product.ID)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if found.Name != "iPhone" {
		t.Fatalf(
			"expected name %q, got %q",
			"iPhone",
			found.Name,
		)
	}
}

func TestProductServiceGetProductByIDInvalidID(t *testing.T) {
	repo := newMockProductRepository()
	cache := newMockProductCache()

	service := NewProductService(
		repo,
		cache,
	)
	_, err := service.GetProductByID(0)

	if err == nil {
		t.Fatal("expected error for ID 0")
	}
}

func TestProductServiceUpdateProduct(t *testing.T) {
	repo := newMockProductRepository()
	cache := newMockProductCache()

	service := NewProductService(
		repo,
		cache,
	)
	product := &model.Product{
		Name:  "Old Product",
		Price: 1000,
		Stock: 5,
	}

	if err := service.CreateProduct(product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	product.Name = "Updated Product"
	product.Price = 1500
	product.Stock = 10

	err := service.UpdateProduct(product)
	if err != nil {
		t.Fatalf("failed to update product: %v", err)
	}

	updated, err := service.GetProductByID(product.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated product: %v", err)
	}

	if updated.Name != "Updated Product" {
		t.Fatalf(
			"expected name %q, got %q",
			"Updated Product",
			updated.Name,
		)
	}

	if updated.Price != 1500 {
		t.Fatalf(
			"expected price %.2f, got %.2f",
			1500.0,
			updated.Price,
		)
	}
}

func TestProductServiceUpdateProductNotFound(t *testing.T) {
	repo := newMockProductRepository()
	cache := newMockProductCache()

	service := NewProductService(
		repo,
		cache,
	)
	product := &model.Product{
		ID:    999,
		Name:  "Missing Product",
		Price: 1000,
		Stock: 5,
	}

	err := service.UpdateProduct(product)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(
			"expected record not found error, got %v",
			err,
		)
	}
}

func TestProductServiceDeleteProduct(t *testing.T) {
	repo := newMockProductRepository()
	cache := newMockProductCache()

	service := NewProductService(
		repo,
		cache,
	)
	product := &model.Product{
		Name:  "Delete Me",
		Price: 1000,
		Stock: 5,
	}

	if err := service.CreateProduct(product); err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	err := service.DeleteProduct(product.ID)
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}

	_, err = service.GetProductByID(product.ID)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf(
			"expected record not found after delete, got %v",
			err,
		)
	}
}

func TestProductService_GetProductByID_RedisGetFailureFallsBackToRepository(t *testing.T) {
	repo := newMockProductRepository()

	expectedProduct := &model.Product{
		ID:          1,
		Name:        "MacBook Pro",
		Description: "Apple laptop",
		Price:       150000,
		Stock:       10,
	}

	repo.products[1] = expectedProduct

	cache := newMockProductCache()
	cache.getErr = errors.New("redis connection failed")

	productService := NewProductService(repo, cache)

	product, err := productService.GetProductByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}

	if product.ID != expectedProduct.ID {
		t.Fatalf("expected product ID %d, got %d",
			expectedProduct.ID,
			product.ID,
		)
	}

	if product.Name != expectedProduct.Name {
		t.Fatalf("expected product name %q, got %q",
			expectedProduct.Name,
			product.Name,
		)
	}
}

func TestProductService_GetProductByID_RedisSetFailureStillReturnsProduct(t *testing.T) {
	repo := newMockProductRepository()

	expectedProduct := &model.Product{
		ID:          2,
		Name:        "Dell Monitor",
		Description: "4K Monitor",
		Price:       50000,
		Stock:       15,
	}

	repo.products[2] = expectedProduct

	cache := newMockProductCache()
	cache.setErr = errors.New("redis set failed")

	productService := NewProductService(repo, cache)

	product, err := productService.GetProductByID(2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}

	if product.ID != expectedProduct.ID {
		t.Fatalf("expected product ID %d, got %d",
			expectedProduct.ID,
			product.ID,
		)
	}
}

func TestProductService_UpdateProduct_RedisDeleteFailureDoesNotFailUpdate(t *testing.T) {
	repo := newMockProductRepository()

	repo.products[3] = &model.Product{
		ID:          3,
		Name:        "Old Product",
		Description: "Old Description",
		Price:       1000,
		Stock:       10,
	}

	cache := newMockProductCache()
	cache.deleteErr = errors.New("redis delete failed")

	productService := NewProductService(repo, cache)

	updatedProduct := &model.Product{
		ID:          3,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       2000,
		Stock:       20,
	}

	err := productService.UpdateProduct(updatedProduct)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	product := repo.products[3]

	if product.Name != "Updated Product" {
		t.Fatalf("expected updated product name, got %q", product.Name)
	}

	if product.Price != 2000 {
		t.Fatalf("expected updated price 2000, got %v", product.Price)
	}
}

func TestProductService_DeleteProduct_RedisDeleteFailureDoesNotFailDelete(t *testing.T) {
	repo := newMockProductRepository()

	repo.products[4] = &model.Product{
		ID:    4,
		Name:  "Delete Product",
		Price: 5000,
		Stock: 5,
	}

	cache := newMockProductCache()
	cache.deleteErr = errors.New("redis delete failed")

	productService := NewProductService(repo, cache)

	err := productService.DeleteProduct(4)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := repo.products[4]; exists {
		t.Fatal("expected product to be deleted from repository")
	}
}
