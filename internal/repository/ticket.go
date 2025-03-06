package repository

import (
	"context"
	"errors"

	"github.com/rickyroynardson/booking/internal/entity"
	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{
		DB: db,
	}
}

func (r *TicketRepository) FindById(ctx context.Context, id string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.DB.WithContext(ctx).Model(&entity.Ticket{}).Where("tickets.id = ?", id).First(&ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrTicketNotFound
		}
		return nil, err
	}
	return &ticket, nil
}
