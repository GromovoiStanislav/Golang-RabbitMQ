package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)
	

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost"
	}

	// initiating the connection with rabbitmq
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("%s: %s", "Can't connect to AMQP", err)
	}
	defer conn.Close()

	// Create a channel
	/*
		Channel opens a unique, concurrent server channel to process the bulk of AMQP
		messages.  Any error from methods on this receiver will render the receiver
		invalid and a new Channel should be opened.
	*/
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Can't create a amqpChannel", err)
	}
	defer amqpChannel.Close()

	// delare the que
	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Could not declare `add` queue", err)
	}

	limit := 10
	var wg sync.WaitGroup
	wg.Add(limit)

	for i := 0; i < limit; i++ {
		number := i
		go func() {
			defer wg.Done()
			err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(fmt.Sprintf("Data %v", number)),
			})

			if err != nil {
				log.Fatalf("Error publishing message: %s", err)
			}
		}()
	}

	wg.Wait()

	fmt.Println("done")
}
