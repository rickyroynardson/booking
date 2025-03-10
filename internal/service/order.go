package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rickyroynardson/booking/internal/entity"
	"github.com/rickyroynardson/booking/internal/messaging/publisher"
	"github.com/rickyroynardson/booking/internal/repository"
)

type OrderService struct {
	repository       *repository.OrderRepository
	ticketRepository *repository.TicketRepository
	bookingPublisher *publisher.BookingPublisher
}

func NewOrderService(repository *repository.OrderRepository, ticketRepository *repository.TicketRepository, bookingPublisher *publisher.BookingPublisher) *OrderService {
	return &OrderService{
		repository,
		ticketRepository,
		bookingPublisher,
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

	orderID := uuid.New().String()
	order := entity.Order{
		ID:       orderID,
		TicketID: body.TicketID,
		Status:   entity.Created,
		Quantity: body.Quantity,
	}
	err = s.repository.Book(ctx, order)
	if err != nil {
		return err
	}

	err = s.bookingPublisher.PublishBooking(orderID, body.TicketID, body.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) Reserve(ctx context.Context, body entity.ReserveOrderRequest) error {
	order, err := s.repository.FindById(ctx, body.OrderID)
	if err != nil {
		return err
	}
	if order.Status != entity.Created {
		return errors.New("order status is not created")
	}

	err = s.repository.Reserve(ctx, body)
	if err != nil {
		return err
	}

	return nil
}
