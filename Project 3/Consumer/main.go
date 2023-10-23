package main

import (
	"fmt"
	"log"
	"os"

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
		"hello", // Имя очереди
		false,   // Долговечность
		false,   // Автоматическое удаление при завершении
		false,   // Исключительность
		false,   // Непосредственная доставка
		nil,     // Аргументы
	)
	failOnError(err, "Failed to declare a queue")

	// Регистрация консьюмера
	msgs, err := ch.Consume(
		q.Name, // Имя очереди
		"",     // Имя консьюмера (пустая строка для автоматической генерации)
		true,   // Автоподтверждение сообщения
		false,  // Эксклюзивность
		false,  // noLocal
		false,  // Не ждать подтверждения от других консьюмеров
		nil,    // Аргументы
	)
	failOnError(err, "Failed to register a consumer")

	// Ожидание и обработка сообщений
	//forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	//<-forever

	// или
	select {}
}
