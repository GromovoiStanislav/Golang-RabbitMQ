package main

import (
	"log"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost"
	}

	connection, _ := amqp.Dial(url)
	defer connection.Close()

	go func(con *amqp.Connection) {
		channel, _ := connection.Channel()
		defer channel.Close()

		durable, exclusive := false, false
		autoDelete, noWait := true, true
		q, _ := channel.QueueDeclare("test", durable, autoDelete, exclusive, noWait, nil)
		channel.QueueBind(q.Name, "#", "amq.topic", false, nil)
		autoAck, exclusive, noLocal, noWait := false, false, false, false
		messages, _ := channel.Consume(q.Name, "", autoAck, exclusive, noLocal, noWait, nil)
		multiAck := false
		for msg := range messages {
			fmt.Println("Body:", string(msg.Body), "Timestamp:", msg.Timestamp)
			msg.Ack(multiAck)
		}
	}(connection)

	go func(con *amqp.Connection) {
		timer := time.NewTicker(1 * time.Second)
		channel, _ := connection.Channel()

		for t := range timer.C {
			msg := amqp.Publishing{
				DeliveryMode: 1,
				Timestamp:    t,
				ContentType:  "text/plain",
				Body:         []byte("Hello world"),
			}
			mandatory, immediate := false, false
			channel.Publish("amq.topic", "ping", mandatory, immediate, msg)
		}
	}(connection)

	select {}
}
