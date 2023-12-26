package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://localhost:5672"
	}

	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Объявление очереди для отложенных сообщений (Dead Letter Queue)
	deadLetterQueue := "dead_letter_queue"
	_, err = ch.QueueDeclare(deadLetterQueue, false, false, false, false, nil)
	failOnError(err, "Failed to declare a dead letter queue")

	// Объявление обмена для отложенных сообщений
	deadLetterExchange := "dead_letter_exchange"
	err = ch.ExchangeDeclare(deadLetterExchange, "direct", false, false, false, false, nil)
	failOnError(err, "Failed to declare a dead letter exchange")

	// Привязка отложенной очереди к обмену
	err = ch.QueueBind(deadLetterQueue, "dead_letter_routing_key", deadLetterExchange, false, nil)
	failOnError(err, "Failed to bind the dead letter queue to the exchange")

	log.Println("Queues and exchanges created successfully.")

	// Объявление основной очереди с TTL
	mainQueue := "main_queue"
	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = deadLetterExchange
	args["x-dead-letter-routing-key"] = "dead_letter_routing_key"
	args["x-message-ttl"] = int32(1000 * 60) // Установка TTL сообщений по умолчанию в миллисекундах (60 секунд)
	_, err = ch.QueueDeclare(mainQueue, false, false, false, false, args)
	failOnError(err, "Failed to declare a main queue")


	message := "Hello, RabbitMQ!"


	// Установка TTL для конкретного сообщения
	expirationTimeInMillis := int32(1000 * 20) // Установка TTL конкретного сообщения в миллисекундах (20 секунд)
	err = ch.Publish("", mainQueue, false, false, amqp.Publishing{
		Body:      []byte(message),
		Expiration: strconv.Itoa(int(expirationTimeInMillis)),
	})
	failOnError(err, "Failed to publish a message with TTL")

	log.Printf("Sent message: %s with TTL %d ms\n", message, expirationTimeInMillis)


	// TTL по умолчанию
	err = ch.Publish("", mainQueue, false, false, amqp.Publishing{
		Body: []byte(message),
	})
	failOnError(err, "Failed to publish a message with common TTL")

	
	log.Printf("Sent message: %s with common TTL\n", message)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
