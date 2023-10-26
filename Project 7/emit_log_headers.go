package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	const exchangeName = "logs_header"

	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	failOnError(err, "Error loading .env file:")

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost:5672"
	}

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"headers",      // type
		false,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// args := os.Args[1:]
	// msg := "Hello!!!"
	// if len(args) > 0 {
	// 	msg = args[0]
	//}

	publishMessage := func(account, method, message string) {
		err := ch.Publish(
			exchangeName,
			"",
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{"account": account, "method": method},
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
		failOnError(err, "Failed to publish a message")

		log.Printf("Sent: %s %s\n", account, method)
	}

	publishMessage("old", "github", "old github") // or msg
	publishMessage("new", "github", "new github") // or msg
	publishMessage("new", "google", "new google") // or msg
	publishMessage("old", "facebook", "old facebook") // or msg

	log.Println("Messages sent")
}
