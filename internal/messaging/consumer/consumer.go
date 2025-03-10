package consumer

import amqp "github.com/rabbitmq/amqp091-go"

type BaseConsumer struct {
	ch    *amqp.Channel
	queue amqp.Queue
	done  chan struct{}
}

func NewBaseConsumer(ch *amqp.Channel, queue amqp.Queue) *BaseConsumer {
	return &BaseConsumer{
		ch,
		queue,
		make(chan struct{}),
	}
}

func (c *BaseConsumer) Stop() error {
	close(c.done)
	return nil
}
