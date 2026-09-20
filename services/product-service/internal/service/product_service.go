package service

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/cache"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/repository"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductByID(id uint) (*model.Product, error)
	GetProducts() ([]model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id uint) error
}

type productService struct {
	repo  repository.ProductRepository
	cache cache.ProductCache
}

func NewProductService(
	repo repository.ProductRepository,
	productCache cache.ProductCache,
) ProductService {
	return &productService{
		repo:  repo,
		cache: productCache,
	}
}

func (s *productService) CreateProduct(product *model.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}

	if product.Price <= 0 {
		return errors.New("product price must be greater than zero")
	}

	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	if err := s.repo.Create(product); err != nil {
		return err
	}

	// A newly created product cannot normally have a stale cache entry.
	// Cache invalidation is best-effort.
	ctx := context.Background()
	if err := s.cache.DeleteProduct(ctx, product.ID); err != nil {
		log.Printf("warning: failed to invalidate product cache for product %d: %v", product.ID, err)
	}

	return nil
}

func (s *productService) GetProductByID(id uint) (*model.Product, error) {
	if id == 0 {
		return nil, errors.New("product ID must be greater than zero")
	}

	ctx := context.Background()

	// Redis is a best-effort cache.
	// If Redis fails, continue to PostgreSQL.
	cachedProduct, err := s.cache.GetProduct(ctx, id)
	if err != nil {
		log.Printf("warning: Redis GET failed for product %d: %v", id, err)
	} else if cachedProduct != nil {
		return cachedProduct, nil
	}

	// PostgreSQL remains the source of truth.
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache population is best-effort.
	if err := s.cache.SetProduct(ctx, product); err != nil {
		log.Printf("warning: Redis SET failed for product %d: %v", id, err)
	}

	return product, nil
}

func (s *productService) GetProducts() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) UpdateProduct(product *model.Product) error {
	if product.ID == 0 {
		return errors.New("product ID is required")
	}

	if product.Name == "" {
		return errors.New("product name is required")
	}

	if product.Price <= 0 {
		return errors.New("product price must be greater than zero")
	}

	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	_, err := s.repo.GetByID(product.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}

		return err
	}

	if err := s.repo.Update(product); err != nil {
		return err
	}

	// Database update succeeded.
	// Cache invalidation is best-effort.
	ctx := context.Background()
	if err := s.cache.DeleteProduct(ctx, product.ID); err != nil {
		log.Printf("warning: failed to invalidate product cache for product %d: %v", product.ID, err)
	}

	return nil
}

func (s *productService) DeleteProduct(id uint) error {
	if id == 0 {
		return errors.New("product ID must be greater than zero")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Database deletion succeeded.
	// Cache invalidation is best-effort.
	ctx := context.Background()
	if err := s.cache.DeleteProduct(ctx, id); err != nil {
		log.Printf("warning: failed to invalidate product cache for product %d: %v", id, err)
	}

	return nil
}
