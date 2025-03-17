package repository

import (
	"context"
	"errors"

	"github.com/rickyroynardson/booking/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func (r *OrderRepository) FindById(ctx context.Context, id string) (*entity.Order, error) {
	var order entity.Order
	if err := r.DB.WithContext(ctx).Model(&entity.Order{}).Where("orders.id = ?", id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Book(ctx context.Context, body entity.Order) error {
	tracer := otel.Tracer("order-repository")
	ctx, span := tracer.Start(ctx, "Book")
	defer span.End()

	span.SetAttributes(
		attribute.String("order_id", body.ID),
		attribute.String("ticket_id", body.TicketID),
		attribute.Int("quantity", body.Quantity),
	)

	err := r.DB.WithContext(ctx).Create(&body).Error
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "order created successfully")
	return nil
}

func (r *OrderRepository) Reserve(ctx context.Context, body entity.ReserveOrderRequest) error {
	var ticket entity.Ticket
	if err := r.DB.WithContext(ctx).Model(&entity.Ticket{}).Where("tickets.id = ?", body.TicketID).First(&ticket).Error; err != nil {
		if err := r.DB.WithContext(ctx).Model(&entity.Order{}).Where("orders.id = ?", body.OrderID).Update("status", entity.Rejected).Error; err != nil {
			return err
		}
		return errors.New("ticket not found")
	}

	if ticket.RemainingCapacity < body.Quantity {
		if err := r.DB.WithContext(ctx).Model(&entity.Order{}).Where("orders.id = ?", body.OrderID).Update("status", entity.Rejected).Error; err != nil {
			return err
		}
		return errors.New("remaining capacity is not enough")
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Model(&entity.Order{}).Where("orders.id = ?", body.OrderID).Update("status", entity.Reserved).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Model(&entity.Ticket{}).Where("tickets.id = ?", body.TicketID).Update("remaining_capacity", gorm.Expr("remaining_capacity - ?", body.Quantity)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		_ = r.DB.WithContext(ctx).Model(&entity.Order{}).Where("orders.id = ?", body.OrderID).Update("status", entity.Rejected).Error
		return err
	}

	return nil
}
