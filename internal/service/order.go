package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rickyroynardson/booking/internal/entity"
	"github.com/rickyroynardson/booking/internal/messaging/publisher"
	"github.com/rickyroynardson/booking/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	tracer := otel.Tracer("order-service")

	ctx, span := tracer.Start(ctx, "Book")
	defer span.End()

	span.SetAttributes(attribute.String("ticket_id", body.TicketID), attribute.Int("quantity", body.Quantity))

	ctx, validationSpan := tracer.Start(ctx, "validate_ticket")
	validationSpan.SetAttributes(attribute.String("ticket_id", body.TicketID))
	ticket, err := s.ticketRepository.FindById(ctx, body.TicketID)
	if err != nil {
		validationSpan.RecordError(err)
		validationSpan.SetStatus(codes.Error, err.Error())
		return err
	}
	validationSpan.End()

	if ticket.RemainingCapacity < body.Quantity {
		span.RecordError(errors.New("remaining ticket is not enough"))
		span.SetStatus(codes.Error, "remaining ticket is not enough")
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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = s.bookingPublisher.PublishBooking(orderID, body.TicketID, body.Quantity)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetAttributes(attribute.String("order_id", orderID))
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
