package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rickyroynardson/booking/config"
	"github.com/rickyroynardson/booking/internal/handler"
	"github.com/rickyroynardson/booking/internal/messaging/consumer"
	"github.com/rickyroynardson/booking/internal/messaging/publisher"
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

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bookings",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare bookings queue: %v", err)
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	bookingPublisher := publisher.NewBookingPublisher(ch, q.Name)

	showRepository := repository.NewShowRepository(db)
	showService := service.NewShowService(showRepository)
	showHandler := handler.NewShowHandler(showService, validator)

	ticketRepository := repository.NewTicketRepository(db)

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository, ticketRepository, bookingPublisher)
	orderHandler := handler.NewOrderHandler(orderService, validator)

	bookingConsumer := consumer.NewBookingConsumer(ch, q, validator, orderService)
	go func() {
		if err := bookingConsumer.Start(); err != nil {
			log.Fatalf("failed to start booking consumer: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "booking")
	})
	e.GET("/shows", showHandler.FindAll)
	e.GET("/shows/:id", showHandler.FindById)
	e.POST("/shows", showHandler.Create)

	e.POST("/shows/:id/book", orderHandler.Book)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get().App.Port)))
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
