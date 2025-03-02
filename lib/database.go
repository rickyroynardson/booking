package lib

import (
	"fmt"

	"github.com/rickyroynardson/booking/config"
	"github.com/rickyroynardson/booking/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Get().DB.Host, config.Get().DB.User, config.Get().DB.Password, config.Get().DB.DBName, config.Get().DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entity.Show{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
