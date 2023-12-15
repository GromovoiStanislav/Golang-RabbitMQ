package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpURL := getAMQPURL()

	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queueName := "my-stream-queue"
	msg := fmt.Sprintf("Hello World! %d", time.Now().Unix())

	err = declareStreamQueue(ch, queueName)
	failOnError(err, "Failed to declare a stream queue")

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	failOnError(err, "Failed to publish a message")

	fmt.Printf(" [x] Sent '%s'\n", msg)
}

func getAMQPURL() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://localhost:5672"
	}

	return amqpURL
}

func declareStreamQueue(ch *amqp.Channel, queueName string) error {
	
	_, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		amqp.Table{ // queue arguments
			amqp.QueueTypeArg:                 amqp.QueueTypeStream,
			amqp.StreamMaxLenBytesArg:         int64(2_000_000_000), // 2 Gb
			//amqp.StreamMaxSegmentSizeBytesArg: int64(500_000_000),          // 500 Mb
			amqp.StreamMaxAgeArg:              "3D",                 // 3 days
			amqp.QueueTTLArg: int32(24 * 60),   // время жизни очереди в минутах
			amqp.QueueMessageTTLArg: int32(60 * 5),   // время жизни сообщений в минутах

		},
	)
	return err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
