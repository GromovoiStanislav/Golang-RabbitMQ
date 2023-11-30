package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)
	

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost"
	}

	// Строка подключения к RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	// Создание очереди
	q, err:= ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	fmt.Println(q)
	if err != nil {
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	// Ждем сообщения
	msgs, err := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println("[*] - Waiting for messages")
	<-forever
}
