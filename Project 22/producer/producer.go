package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProducerConfig struct {
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

type Producer struct {
	config ProducerConfig
	Rabbit *Rabbit
}

func NewProducer(config ProducerConfig, rabbit *Rabbit) *Producer {
	return &Producer{
		config: config,
		Rabbit: rabbit,
	}
}

func (p *Producer) SendMessage(message string) error {
	con, err := p.Rabbit.Connection()
	if err != nil {
		return err
	}

	chn, err := con.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()

	err = chn.ExchangeDeclare(
		p.config.ExchangeName,
		p.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = chn.Publish(
		p.config.ExchangeName,
		p.config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	log.Println("Message sent:", message)
	return nil
}
