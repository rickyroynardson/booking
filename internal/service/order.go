package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rickyroynardson/booking/internal/entity"
	"github.com/rickyroynardson/booking/internal/repository"
)

type OrderService struct {
	repository       *repository.OrderRepository
	ticketRepository *repository.TicketRepository
}

func NewOrderService(repository *repository.OrderRepository, ticketRepository *repository.TicketRepository) *OrderService {
	return &OrderService{
		repository,
		ticketRepository,
	}
}

func (s *OrderService) Book(ctx context.Context, body entity.BookOrderRequest) error {
	ticket, err := s.ticketRepository.FindById(ctx, body.TicketID)
	if err != nil {
		return err
	}
	if ticket.RemainingCapacity < body.Quantity {
		return errors.New("remaining ticket is not enough")
	}

	order := entity.Order{
		ID:       uuid.New().String(),
		TicketID: body.TicketID,
		Status:   entity.Created,
		Quantity: body.Quantity,
	}
	return s.repository.Book(ctx, order)
}
