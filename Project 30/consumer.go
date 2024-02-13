package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}
	
	conn, err := amqp.DialConfig(amqpURL, amqp.Config{
		Properties: amqp.NewConnectionProperties(),
	})
	if err != nil {
		panic(err)
	}

	var channel *amqp.Channel
	if channel, err = conn.Channel(); err != nil {
		panic(err)
	}

	if err := channel.ExchangeDeclare(
		"user_dlx",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		panic(err)
	}
	 
	if _, err := channel.QueueDeclare(
		"user_create_dlx",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		panic(err)
	}
	 
	if err := channel.QueueBind(
		"user_create_dlx",
		"",
		"user_dlx",
		false,
		nil,
	); err != nil {
		panic(err)
	}


	if _, err := channel.QueueDeclare(
		"user_create",
		true,
		false,
		false,
		false,
		amqp.Table{"x-dead-letter-exchange": "user_dlx"},
	); err != nil {
		panic(err)
	}


	go func() {
		messages, _ := channel.Consume("user_create_dlx", "dlx-consumer", false, false, false, false, nil)
		for msg := range messages {
			fmt.Println("dlx", string(msg.Body))
			// msg.Ack(false)	
		}
	}()


	messages, _ := channel.Consume("user_create", "my-consumer", false, false, false, false, nil)
	for msg := range messages {
		go func(msg amqp.Delivery) {
			fmt.Println(string(msg.Body))

			// it will take the message out of the original queue and put into DLX queue
			msg.Nack(false, false)	
		}(msg)
	}


}

