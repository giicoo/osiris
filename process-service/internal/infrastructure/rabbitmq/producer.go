package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (p *AlertProducing) PublicMessage(body []byte) error {
	err := p.amqpChan.Publish(
		p.cfg.Rabbitmq.Exchange.Name,
		p.cfg.Rabbitmq.Notifications.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish rabbitmq: %w", err)
	}
	return nil
}
