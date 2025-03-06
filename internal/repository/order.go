package repository

import (
	"context"

	"github.com/rickyroynardson/booking/internal/entity"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

func (r *OrderRepository) Book(ctx context.Context, body entity.Order) error {
	return r.DB.WithContext(ctx).Create(&body).Error
}
