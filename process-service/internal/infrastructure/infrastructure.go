package infrastructure

import amqp "github.com/rabbitmq/amqp091-go"

type AlertProducing interface {
	PublicMessage(body []byte) error
	ConsumeMessage() (<-chan amqp.Delivery, error)
}
