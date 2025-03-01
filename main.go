package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rickyroynardson/booking/config"
	"github.com/rickyroynardson/booking/internal/repository"
	"github.com/rickyroynardson/booking/lib"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	db, err := lib.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	_ = repository.NewShowRepository(db)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "booking")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get().App.Port)))
}
