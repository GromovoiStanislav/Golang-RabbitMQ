package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	failOnError(err, "Error loading .env file:")

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost"
}

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	confirms := make(chan amqp.Confirmation)
	ch.NotifyPublish(confirms)
	go func() {
		for confirm := range confirms {
			if confirm.Ack {
				// code when messages is confirmed
				log.Printf("Confirmed")
			} else {
				// code when messages is nack-ed
				log.Printf("Nacked")
			}
		}
	}()

	err = ch.Confirm(false)
	failOnError(err, "Failed to confirm")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	consume(ch, q.Name)
	publish(ch, q.Name, "hello")

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	
	var forever chan struct{}
	<-forever
}

func consume(ch *amqp.Channel, qName string) {
	msgs, err := ch.Consume(
		qName, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}

func publish(ch *amqp.Channel, qName, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := ch.PublishWithContext(ctx, "", qName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(text),
	})
	failOnError(err, "Failed to publish a message")
}