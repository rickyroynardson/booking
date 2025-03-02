package entity

type Show struct {
	ID          string `gorm:"type:uuid;not null;primaryKey"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
}

type CreateShowRequest struct {
	Name        string `validate:"required,min=6,max=255"`
	Description string
}
