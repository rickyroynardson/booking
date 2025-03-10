package entity

import "errors"

type OrderStatus string

const (
	Created  OrderStatus = "created"
	Reserved OrderStatus = "reserved"
	Rejected OrderStatus = "rejected"
)

var ErrOrderNotFound = errors.New("order not found")

type Order struct {
	ID       string `gorm:"type:uuid;not null;primaryKey"`
	TicketID string `gorm:"type:uuid;not null"`
	Status   OrderStatus
	Quantity int `gorm:"type:integer;not null"`
}

type BookOrderRequest struct {
	TicketID string `json:"ticket_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

type ReserveOrderRequest struct {
	OrderID  string `json:"order_id"`
	TicketID string `json:"ticket_id"`
	Quantity int    `json:"quantity"`
}
