package lib

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rickyroynardson/booking/config"
	"github.com/rickyroynardson/booking/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Get().DB.Host, config.Get().DB.User, config.Get().DB.Password, config.Get().DB.DBName, config.Get().DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entity.Show{}, &entity.Ticket{}, &entity.Order{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
