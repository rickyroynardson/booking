package publisher

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type BasePublisher struct {
	ch         *amqp.Channel
	exchange   string
	routingKey string
}

func NewBasePublisher(ch *amqp.Channel, exchange, routingKey string) *BasePublisher {
	return &BasePublisher{
		ch,
		exchange,
		routingKey,
	}
}

func (p *BasePublisher) Publish(message []byte) error {
	return p.ch.Publish(
		p.exchange,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
