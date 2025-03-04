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

func (s *ShowService) FindAll(ctx context.Context, req entity.FindAllShowRequest) (*entity.PaginationResponse, error) {
	shows, totalRecords, err := s.repository.FindAll(ctx, req.Search, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(totalRecords) / req.Limit
	if int(totalRecords)%req.Limit > 0 {
		totalPages++
	}

	metadata := entity.PaginationMetadata{
		CurrentPage:  req.Page,
		PageSize:     req.Limit,
		TotalRecords: int(totalRecords),
		TotalPages:   totalPages,
		HasNext:      req.Page < totalPages,
		HasPrevious:  req.Page > 1,
	}

	response := &entity.PaginationResponse{
		Data:     shows,
		Metadata: metadata,
	}

	return response, nil
}

func (s *ShowService) Create(ctx context.Context, body entity.CreateShowRequest) error {
	show := entity.Show{
		ID:          uuid.New().String(),
		Name:        body.Name,
		Description: body.Description,
	}
	return s.repository.Create(ctx, show)
}
