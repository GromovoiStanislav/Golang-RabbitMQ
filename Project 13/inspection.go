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

	// Проверьте существование очереди
	
	// Имя очереди, которую вы хотите проверить
	exchangeName := "hello.direct"

	// Попробуйте объявить обменник (если он не существует, вернется ошибка)
	err = ch.ExchangeDeclarePassive(
		exchangeName,
		"fanout", // Тип обменника (fanout, direct, topic, headers)
		false,    // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
	if err != nil {
		if amqpErr, ok := err.(*amqp.Error); ok && amqpErr.Code == amqp.NotFound {
			// Обменник не существует
			fmt.Printf("Обменник %s не существует\n", exchangeName)
		} else {
			// Другая ошибка
			log.Fatalf("Ошибка при проверке обменника: %v", err)
		}
	} else {
		// Обменник существует
		fmt.Printf("Обменник %s существует\n", exchangeName)
	}
	// После выполнения операций с каналом, проверьте, не закрыт ли он	
	if ch.IsClosed() {
		log.Println(ch.IsClosed())
		return
	}


	// Проверьте существование очереди

	// Имя очереди, которую вы хотите проверить
	queueName := "hello"

	// 1-й способ - Deprecated: Use QueueDeclarePassive instead
	_, err = ch.QueueInspect(queueName)
	if err != nil {
		if amqpErr, ok := err.(*amqp.Error); ok && amqpErr.Code == amqp.NotFound {
			// Очередь не существует
			fmt.Printf("Очередь %s не существует\n", queueName)
		} else {
			// Другая ошибка
			log.Printf("Ошибка при проверке очереди: %v", err)
		}
	} else {
		// Очередь существует
		fmt.Printf("Очередь %s существует\n", queueName)
	}
	// После выполнения операций с каналом, проверьте, не закрыт ли он	
	if ch.IsClosed() {
		log.Println(ch.IsClosed())
		return
	}



	// 2-й способ
	_, err = ch.QueueDeclarePassive(
		queueName, // Имя очереди
		false,   // Долговечность
		false,   // Автоматическое удаление при завершении
		false,   // Исключительность
		false,   // Непосредственная доставка
		nil,     // Аргументы
	)
	if err != nil {
		if amqpErr, ok := err.(*amqp.Error); ok {
			if amqpErr.Code == amqp.NotFound {
				// Очередь не существует
				fmt.Printf("Очередь %s не существует\n", queueName)
			} else {
				// Другая ошибка
				log.Printf("Ошибка при проверке очереди: %v", err)
			}
		} else {
			// Не является *amqp.Error, обработка других типов ошибок
			log.Printf("Неожиданный тип ошибки: %v", err)
		}
	}
	// После выполнения операций с каналом, проверьте, не закрыт ли он	
	if ch.IsClosed() {
		log.Println(ch.IsClosed())
		return
	}


	// Отправка сообщения.....

}
