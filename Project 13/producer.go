package main

import (
	//"context"
	"fmt"
	"log"
	"os"
	//"time"

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

	

	// Создание канала для получения подтверждений
	notifyConfirm := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// Установка подтверждений на канале
	err = ch.Confirm(false)
	failOnError(err, "Failed to enable publisher confirms")


	// Отправка сообщения
	body := "Hello, RabbitMQ!"
		

	// Создание контекста с таймаутом, например, 5 секунд
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	//err = ch.PublishWithContext(
	err = ch.Publish(
		//ctx,     // Контекст
		"hello.direct",     // Обмен
		"hello", // Ключ маршрутизации
		false,  // Обязательное
		false,  // Непосредственная доставка
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf("[x] Sent %s\n", body)


	// Ожидание подтверждения асинхронно
	confirmed  := <-notifyConfirm
	if confirmed.Ack {
		fmt.Printf("[x] Message sent successfully: %s\n", body)
	} else {
		fmt.Printf("[!] Failed to send message: %s\n", body)
	}
}
