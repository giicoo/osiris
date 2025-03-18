package infrastructure

import amqp "github.com/rabbitmq/amqp091-go"

type AlertProducing interface {
	ConsumeMessage() (<-chan amqp.Delivery, error)
}
