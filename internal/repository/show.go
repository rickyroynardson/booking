package repository

import (
	"context"

	"github.com/rickyroynardson/booking/internal/entity"
	"gorm.io/gorm"
)

type ShowRepository struct {
	DB *gorm.DB
}

func NewShowRepository(db *gorm.DB) *ShowRepository {
	return &ShowRepository{
		DB: db,
	}
}

func (r *ShowRepository) Create(ctx context.Context, body entity.Show) error {
	return r.DB.WithContext(ctx).Create(&body).Error
}
