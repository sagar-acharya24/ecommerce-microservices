package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/model"
)

type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id uint) (*model.Product, error)
	GetAll() ([]model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product

	result := r.db.First(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (r *productRepository) GetAll() ([]model.Product, error) {
	var products []model.Product

	result := r.db.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *productRepository) Update(product *model.Product) error {
	result := r.db.Model(&model.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]interface{}{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
		})

	return result.Error
}

func (r *productRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
