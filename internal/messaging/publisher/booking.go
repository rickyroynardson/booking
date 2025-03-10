package publisher

import (
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BookingPublisher struct {
	*BasePublisher
}

type BookingMessage struct {
	OrderID  string `json:"order_id"`
	TicketID string `json:"ticket_id"`
	Quantity int    `json:"quantity"`
}

func NewBookingPublisher(ch *amqp.Channel, queueName string) *BookingPublisher {
	return &BookingPublisher{
		BasePublisher: NewBasePublisher(ch, "", queueName),
	}
}

func (p *BookingPublisher) PublishBooking(orderID, ticketID string, quantity int) error {
	message := BookingMessage{
		OrderID:  orderID,
		TicketID: ticketID,
		Quantity: quantity,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// TEMP LOG, DELETE LATER
	f, err := os.OpenFile("order_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening log file: %v", err)
	} else {
		defer f.Close()
		logger := log.New(f, "", log.LstdFlags)
		logger.Printf("order created: %s", orderID)
	}
	//

	return p.Publish(b)
}
