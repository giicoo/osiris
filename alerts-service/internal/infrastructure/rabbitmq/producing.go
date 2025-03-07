package rabbitmq

import (
	"fmt"

	"github.com/giicoo/osiris/alerts-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type AlertProducing struct {
	conn     *amqp.Connection
	amqpChan *amqp.Channel
	cfg      *config.Config
}

func NewAlertProducing(cfg *config.Config) *AlertProducing {
	return &AlertProducing{
		cfg: cfg,
	}
}

func (p *AlertProducing) InitAlertProducing() error {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return fmt.Errorf("connect rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel rabbitmq: %w", err)
	}
	p.amqpChan = ch
	p.conn = conn

	if err := p.Setup(); err != nil {
		return fmt.Errorf("setup rabbitmq: %w", err)
	}
	return nil
}

func (p *AlertProducing) Setup() error {
	err := p.amqpChan.ExchangeDeclare(
		p.cfg.Rabbitmq.Exchange.Name,
		p.cfg.Rabbitmq.Exchange.Type,
		p.cfg.Rabbitmq.Exchange.Durability,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("setup  rabbitmq: %w", err)
	}
	queue, err := p.amqpChan.QueueDeclare(
		p.cfg.Rabbitmq.Queue.Name,
		p.cfg.Rabbitmq.Queue.Durability,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("setup rabbitmq: %w", err)
	}

	err = p.amqpChan.QueueBind(
		queue.Name,
		"",
		p.cfg.Rabbitmq.Exchange.Name,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("setup rabbitmq: %w", err)
	}
	return nil
}

func (p *AlertProducing) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("close connection: %w", err)
	}
	if err := p.amqpChan.Close(); err != nil {
		return fmt.Errorf("close chan: %w", err)
	}
	return nil
}

func (p *AlertProducing) PublicMessage(body []byte) error {
	err := p.amqpChan.Publish(
		p.cfg.Rabbitmq.Exchange.Name,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "json",
			Body:        body,
		},
	)
	logrus.Info(err)
	if err != nil {
		return fmt.Errorf("publish rabbitmq: %w", err)
	}
	return nil
}
