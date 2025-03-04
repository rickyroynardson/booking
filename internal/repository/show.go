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

func (r *ShowRepository) FindAll(ctx context.Context, search string, page, limit int) ([]entity.Show, int64, error) {
	var shows []entity.Show
	var totalRecords int64

	query := r.DB.WithContext(ctx)
	if search != "" {
		query = query.Where("name ILIKE ?", search)
	}

	if err := query.Model(&entity.Show{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&shows).Error; err != nil {
		return nil, 0, err
	}
	return shows, totalRecords, nil
}

func (r *ShowRepository) Create(ctx context.Context, body entity.Show) error {
	return r.DB.WithContext(ctx).Create(&body).Error
}
