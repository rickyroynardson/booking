package repository

import (
	"context"
	"errors"

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

func (r *ShowRepository) FindAll(ctx context.Context, search string, page, limit int) ([]entity.FindAllShowResponse, int64, error) {
	var shows []entity.FindAllShowResponse
	var totalRecords int64

	query := r.DB.WithContext(ctx).Model(&entity.Show{})
	if search != "" {
		query = query.Where("shows.name ILIKE ?", search)
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Select("shows.id, shows.name, shows.description, COUNT(tickets.id) as tickets").Joins("LEFT JOIN tickets ON shows.id = tickets.show_id").Group("shows.id").Find(&shows).Error; err != nil {
		return nil, 0, err
	}
	return shows, totalRecords, nil
}

func (r *ShowRepository) FindById(ctx context.Context, id string) (*entity.FindShowByIdResponse, error) {
	var show entity.FindShowByIdResponse

	if err := r.DB.WithContext(ctx).Model(&entity.Show{}).Where("shows.id = ?", id).Preload("Tickets").First(&show).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrShowNotFound
		}
		return nil, err
	}

	return &show, nil
}

func (r *ShowRepository) Create(ctx context.Context, body entity.Show) error {
	return r.DB.WithContext(ctx).Create(&body).Error
}
