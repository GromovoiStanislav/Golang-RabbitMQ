package main

import (
	"log"
	"os"
	"time"

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

	mainExchange := "main_exchange"
	mainQueue := "main_queue"
	deadLetterExchange := "dead_letter_exchange"
	deadLetterQueue := "dead_letter_queue"
	rightRoutingKey := "rightkey"
	invalidRoutingKey := "invalidKey"

	// Объявление обмена для основных сообщений
	err = ch.ExchangeDeclare(mainExchange, "direct", false, false, false, false, amqp.Table{"alternate-exchange": deadLetterExchange})
	failOnError(err, "Failed to declare main exchange")

	// Объявление основной очереди
	_, err = ch.QueueDeclare(mainQueue, false, false, false, false, nil)
	failOnError(err, "Failed to declare main queue")

	// Привязка "main_queue" к обмену для отложенных сообщений
	err = ch.QueueBind(mainQueue, rightRoutingKey, mainExchange, false, nil)
	failOnError(err, "Failed to bind main queue to main exchange with rightRoutingKey")

	// Объявление обмена для отложенных сообщений
	err = ch.ExchangeDeclare(deadLetterExchange, "fanout", false, false, false, false, nil)
	failOnError(err, "Failed to declare dead letter exchange")

	// Объявление "Dead Letter Queue"
	_, err = ch.QueueDeclare(deadLetterQueue, false, false, false, false, nil)
	failOnError(err, "Failed to declare dead letter queue")

	// Привязка "Dead Letter Queue" к обмену для отложенных сообщений
	err = ch.QueueBind(deadLetterQueue, "", deadLetterExchange, false, nil)
	failOnError(err, "Failed to bind dead letter queue to dead letter exchange with # routing key")

	log.Println("Exchanges and queues created successfully.")

	message := "Hello, World!"

	// Отправка сообщения с правильным ключом маршрутизации
	err = ch.Publish(mainExchange, rightRoutingKey, false, false, amqp.Publishing{
		Body: []byte(message),
	})
	failOnError(err, "Failed to publish a message with right routing key")

	log.Printf("Sent message: %s with routing key: %s\n", message, rightRoutingKey)

	// Отправка сообщения с неправильным ключом маршрутизации
	err = ch.Publish(mainExchange, invalidRoutingKey, false, false, amqp.Publishing{
		Body: []byte(message),
	})
	failOnError(err, "Failed to publish a message with invalid routing key")

	log.Printf("Sent message: %s with routing key: %s\n", message, invalidRoutingKey)

	// Закрытие соединения
	time.Sleep(500 * time.Millisecond)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
