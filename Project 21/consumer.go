package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)
	

func main() {
	amqpURL := getAMQPURL()

	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queueName := "my-stream-queue"
	err = declareStreamQueue(ch, queueName)
	failOnError(err, "Failed to declare a stream queue")

	err = ch.Qos(10, 0, false)
	failOnError(err, "Failed to set QoS")

	args := make(amqp.Table)
	args["x-stream-offset"] = "first"// начинать с начала стрима
	// "x-stream-offset": "first"
	// "x-stream-offset": "last"
	// "x-stream-offset": "next"
	// "x-stream-offset": 5
	// "x-stream-offset": { '!': 'timestamp', value: 1692956499 }

	msgs, err := ch.Consume(
		queueName,
		"consumer-name",
		false,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to register a consumer")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Interrupt signal received. Closing channel and connection...")
		ch.Close()
		conn.Close()
		os.Exit(0)
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")

	for msg := range msgs {
		fmt.Printf(" [x] Received '%s'\n", msg.Body)
		msg.Ack(false)
	}
}

func getAMQPURL() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://localhost:5672"
	}

	return amqpURL
}

func declareStreamQueue(ch *amqp.Channel, queueName string) error {
	args := make(amqp.Table)
	args["x-queue-type"] = "stream"
	args["x-max-length-bytes"] = int64(2_000_000_000) // 2 Gb
	//args["x-max-segment-size-bytes"] = int64(500_000_000) // 500 Mb
	args["x-max-age"] = "3D"                 		  // 3 days

	_, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		args,  // arguments
	)
	return err
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}