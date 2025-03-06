package entity

type OrderStatus string

const (
	Created  OrderStatus = "created"
	Reserved OrderStatus = "reserved"
	Rejected OrderStatus = "rejected"
)

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
