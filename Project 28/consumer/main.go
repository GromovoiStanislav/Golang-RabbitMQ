package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = fmt.Sprintf("amqp://%s:%s@%s/%s", "guest", "guest", "localhost:5672", "")
	}
	
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	// Create a new Queue
	_, err = ch.QueueDeclare("events", true, false, false, true, amqp.Table{
		"x-queue-type":                    "stream",
		"x-stream-max-segment-size-bytes": 30000,  // EACH SEGMENT FILE IS ALLOWED 0.03 MB
		"x-max-length-bytes":              150000, // TOTAL STREAM SIZE IS 0.15 MB
	})
	if err != nil {
		panic(err)
	}


	if err := ch.Qos(50, 0, false); err != nil {
		panic(err)
	}

	// Auto ACk has to be FALSE
	stream, err := ch.Consume("events", "events_consumer", false, false, false, false, amqp.Table{
		"x-stream-offset": "10m", // STARTING POINT
	})
	if err != nil {
		panic(err)
	}

	// Loop forever and just read the messages
	fmt.Println("Starting to consume stream")
	for event := range stream {
		fmt.Printf("Event: %s\n", event.CorrelationId)
		fmt.Printf("Headers: %v\n", event.Headers)
		// The payload is in the body
		fmt.Printf("Data: %v\n\n", string(event.Body))

	}
	ch.Close()
}