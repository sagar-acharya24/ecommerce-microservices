package repository

import (
	"errors"

	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	GetByID(id uint) (*model.Order, error)
	GetByUserID(userID uint) ([]model.Order, error)
	Update(order *model.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order

	result := r.db.First(&order, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func (r *orderRepository) GetByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order

	result := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (r *orderRepository) Update(order *model.Order) error {
	result := r.db.
		Model(&model.Order{}).
		Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"status": order.Status,
		})

	return result.Error
}
