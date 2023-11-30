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
			fmt.Println(err)
			panic(err)
		}
	
		// Отправка сообщения
		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("Hello World"),
			},
		)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Successfully Published Message to Queue")
}
