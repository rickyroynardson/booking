package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rickyroynardson/booking/internal/entity"
	"github.com/rickyroynardson/booking/internal/repository"
)

type ShowService struct {
	repository *repository.ShowRepository
}

func NewShowService(repository *repository.ShowRepository) *ShowService {
	return &ShowService{
		repository,
	}
}

func (s *ShowService) Create(ctx context.Context, body entity.CreateShowRequest) error {
	show := entity.Show{
		ID:          uuid.New().String(),
		Name:        body.Name,
		Description: body.Description,
	}
	return s.repository.Create(ctx, show)
}
