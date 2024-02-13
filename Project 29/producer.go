package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()




	message := "Hello, RabbitMQ!"



	err = ch.Publish("my-exchange", "my-routing-key", false, false, amqp.Publishing{
		Body: []byte(message),
		MessageId: "some-unique-id",
	})
	failOnError(err, "Failed to publish a message with common TTL")

	fmt.Printf("Sent message: %s\n", message)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
