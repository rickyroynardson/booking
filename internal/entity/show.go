package entity

import "errors"

var ErrShowNotFound = errors.New("show not found")

type Show struct {
	ID          string   `gorm:"type:uuid;not null;primaryKey"`
	Name        string   `gorm:"type:varchar(255);not null"`
	Description string   `gorm:"type:text"`
	Tickets     []Ticket `gorm:"foreignKey:ShowID"`
}

type FindAllShowRequest struct {
	PaginationRequest
	Search string `query:"search"`
}

type FindAllShowResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tickets     int    `json:"tickets"`
}

type FindShowByIdResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tickets     []Ticket `gorm:"foreignKey:ShowID;references:ID"`
}

type CreateShowRequest struct {
	Name        string `validate:"required,min=6,max=255"`
	Description string
}
