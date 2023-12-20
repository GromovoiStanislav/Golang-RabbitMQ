package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	rabbitConfig := RabbitConfig{
		Schema:         "amqp",
		Username:       os.Getenv("CLOUDAMQP_USERNAME"),
		Password:       os.Getenv("CLOUDAMQP_PASSWORD"),
		Host:           os.Getenv("CLOUDAMQP_HOST"),
		Port:           os.Getenv("CLOUDAMQP_PORT"),
		VHost:          os.Getenv("CLOUDAMQP_VHOST"),
		ConnectionName: "my_app_name",
	}

	rabbit := NewRabbit(rabbitConfig)
	if err := rabbit.Connect(); err != nil {
		log.Fatalln("unable to connect to rabbit", err)
	}

	// Producer
	pc := ProducerConfig{
		ExchangeName: "user",
		ExchangeType: "direct",
		RoutingKey:   "create",
	}
	producer := NewProducer(pc, rabbit)

	message := "Hello, RabbitMQ!"
	if err := producer.SendMessage(message); err != nil {
		log.Println("Error sending message:", err)
	}

	
}
