package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rickyroynardson/booking/config"
	"github.com/rickyroynardson/booking/internal/handler"
	"github.com/rickyroynardson/booking/internal/repository"
	"github.com/rickyroynardson/booking/internal/service"
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

	validator := validator.New(validator.WithRequiredStructEnabled())

	showRepository := repository.NewShowRepository(db)
	showService := service.NewShowService(showRepository)
	showHandler := handler.NewShowHandler(showService, validator)

	ticketRepository := repository.NewTicketRepository(db)

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository, ticketRepository)
	orderHandler := handler.NewOrderHandler(orderService, validator)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "booking")
	})
	e.GET("/shows", showHandler.FindAll)
	e.GET("/shows/:id", showHandler.FindById)
	e.POST("/shows", showHandler.Create)

	e.POST("/shows/:id/book", orderHandler.Book)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get().App.Port)))
}
