package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("RabbitMQ in Golang: Getting started tutorial")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://localhost:5672"
	}

	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	// opening a channel over the connection established to interact with RabbitMQ
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()



	args := make(amqp.Table)
	args["x-message-ttl"] = 60*1000 // Установка TTL сообщений по умолчанию в миллисекундах (10 секунд)
	args["x-expires"] = 2*60*1000 // Установка TTL самой очереди

	queue, err := channel.QueueDeclare(
		"testing", // name
		false,     // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		args,       // args
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Queue status:", queue)


	// publishing a message
	err = channel.Publish(
		"",        // exchange
		"testing", // key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Test Message"),
		},
	)
	if err != nil {
		panic(err)
	}

	// publishing a message
	err = channel.Publish(
		"",        // exchange
		"testing", // key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Test Message 30 sec"),
			Expiration:  "20000",// Установка TTL конкретного сообщения в миллисекундах (20 секунд)
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully published messages")
}
