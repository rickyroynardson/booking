package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rickyroynardson/booking/internal/entity"
	"github.com/rickyroynardson/booking/internal/service"
)

type BookingMessage struct {
	OrderID  string `json:"order_id" validate:"required"`
	TicketID string `json:"ticket_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

type BookingConsumer struct {
	*BaseConsumer
	validator    *validator.Validate
	orderService *service.OrderService
}

func NewBookingConsumer(ch *amqp.Channel, queue amqp.Queue, validator *validator.Validate, orderService *service.OrderService) *BookingConsumer {
	return &BookingConsumer{
		BaseConsumer: NewBaseConsumer(ch, queue),
		validator:    validator,
		orderService: orderService,
	}
}

func (c *BookingConsumer) Start() error {
	msgs, err := c.ch.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("Booking consumer started. Waiting for messages...")

	go func() {
		for {
			select {
			case <-c.done:
				log.Println("Booking consumer stopped")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Booking consumer channel closed")
					return
				}

				err := c.processMessage(msg)
				if err != nil {
					log.Printf("error processing message: %v", err)
					msg.Nack(false, false)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

func (c *BookingConsumer) processMessage(msg amqp.Delivery) error {
	var bookingMsg BookingMessage
	if err := json.Unmarshal(msg.Body, &bookingMsg); err != nil {
		log.Printf("error unmarshaling message: %v", err)
		return err
	}

	if err := c.validator.Struct(bookingMsg); err != nil {
		log.Printf("error validating message: %v", err)
		return err
	}

	err := c.orderService.Reserve(context.Background(), entity.ReserveOrderRequest{
		OrderID:  bookingMsg.OrderID,
		TicketID: bookingMsg.TicketID,
		Quantity: bookingMsg.Quantity,
	})
	if err != nil {
		// TEMP LOG, DELETE LATER
		f, err := os.OpenFile("order_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("error opening log file: %v", err)
		} else {
			defer f.Close()
			logger := log.New(f, "", log.LstdFlags)
			logger.Printf("order failed to process: %s", bookingMsg.OrderID)
		}
		//
		log.Printf("failed to process booking: %v", err)
		return err
	}

	// TEMP LOG, DELETE LATER
	f, err := os.OpenFile("order_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening log file: %v", err)
	} else {
		defer f.Close()
		logger := log.New(f, "", log.LstdFlags)
		logger.Printf("order processed: %s", bookingMsg.OrderID)
	}
	//

	log.Printf("booking processed successfully for ticket: %s", bookingMsg.TicketID)
	return nil
}
