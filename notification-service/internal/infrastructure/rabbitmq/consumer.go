package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (p *AlertProducing) ConsumeMessage() (<-chan amqp.Delivery, error) {
	msgs, err := p.amqpChan.Consume(
		p.cfg.Rabbitmq.Notifications.Name, // очередь
		"",                                // consumer
		true,                              // auto-ack
		false,                             // exclusive
		false,                             // no-local
		false,                             // no-wait
		nil,                               // args
	)
	if err != nil {
		return nil, fmt.Errorf("get chan delivery: %w", err)
	}
	return msgs, nil
}
