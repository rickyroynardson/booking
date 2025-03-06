package entity

import "errors"

var ErrTicketNotFound = errors.New("ticket not found")

type Ticket struct {
	ID                string `gorm:"type:uuid;not null;primaryKey"`
	ShowID            string `gorm:"type:uuid;not null"`
	Name              string `gorm:"type:varchar(255);not null"`
	Price             int    `gorm:"type:integer;not null;default:0"`
	Capacity          int    `gorm:"type:integer;not null;default:0"`
	RemainingCapacity int    `gorm:"type:integer;not null;default:0"`
}
