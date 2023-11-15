package main

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/joho/godotenv"

	"go-rabbitmq-example/lib/event"
)

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672"
	}

	// Строка подключения к RabbitMQ
	connection, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	consumer, err := event.NewConsumer(connection)
	if err != nil {
		panic(err)
	}
	consumer.Listen(os.Args[1:])
}