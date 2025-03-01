package repository

import "gorm.io/gorm"

type ShowRepository struct {
	DB *gorm.DB
}

func NewShowRepository(db *gorm.DB) *ShowRepository {
	return &ShowRepository{
		DB: db,
	}
}
