package main

import (
	"log"
	"os"

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

	// Qos
	// This code sets up Quality of Service (QoS) for the channel.
	// QoS controls how many messages or how many bytes the server will try to keep on the network for consumers before
	// receiving delivery acknowledgments
	err = amqpChannel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("%s: %s", "Could not configure QoS", err)
	}

	autoAck, exclusive, noLocal, noWait := false, false, false, false

	// Consume messages from the queue
	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		autoAck,
		exclusive,
		noLocal,
		noWait,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s", "Could not register consumer", err)
	}

	stopChan := make(chan bool)

	// Start consuming messages
	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	// Stop for program termination
	<-stopChan
}

