package main

import (
	"log"
	"os"
	"time"

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

	consumerConfig := ConsumerConfig{
		ExchangeName:  "user",
		ExchangeType:  "direct",
		RoutingKey:    "create",
		QueueName:     "user_create",
		ConsumerName:  "my_app_name",
		ConsumerCount: 3,
		PrefetchCount: 1,
	}
	consumerConfig.Reconnect.MaxAttempt = 60
	consumerConfig.Reconnect.Interval = 1 * time.Second
	
	consumer := NewConsumer(consumerConfig, rabbit)
	if err := consumer.Start(); err != nil {
		log.Fatalln("unable to start consumer", err)
	}

	select {}
}
