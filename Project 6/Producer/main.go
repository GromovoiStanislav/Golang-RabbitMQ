package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)
	

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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

	// Строка подключения к RabbitMQ
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Создание канала
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Объявление очереди
	q, err := ch.QueueDeclare(
		"my_queue", // Имя очереди
		false,   // Долговечность
		false,   // Автоматическое удаление при завершении
		false,   // Исключительность
		false,   // Непосредственная доставка
		nil,     // Аргументы
	)
	failOnError(err, "Failed to declare a queue")

	// Отправка сообщения в очередь
	body := fmt.Sprintf("Hello, world! (%s)", time.Now())
	headers := make(amqp.Table)
    headers["headerKey"] = "headerValue"
	err = ch.Publish(
		"",     // Обмен
		q.Name, // Ключ маршрутизации
		false,  // Обязательное
		false,  // Непосредственная доставка
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			CorrelationId: "ddd888",
			ReplyTo: "reply_queue",
			Headers: headers,
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf(" [x] Sent %s\n", body)
}
