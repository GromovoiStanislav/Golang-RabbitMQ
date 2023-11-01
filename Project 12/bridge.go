package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)


func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost"
	}


	// This example acts as a bridge, shoveling all messages sent from the source
	// exchange "log1" to destination exchange "log2".

	// Confirming publishes can help from overproduction and ensure every message
	// is delivered.

	// Setup the source of the store and forward
	source, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("connection.open source: %s", err)
	}
	defer source.Close()

	chs, err := source.Channel()
	if err != nil {
		log.Fatalf("channel.open source: %s", err)
	}

	if err := chs.ExchangeDeclare("log1", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("exchange.declare destination: %s", err)
	}

	if _, err := chs.QueueDeclare("page", true, true, false, false, nil); err != nil {
		log.Fatalf("queue.declare source: %s", err)
	}

	if err := chs.QueueBind("page", "alert", "log1", false, nil); err != nil {
		log.Fatalf("queue.bind source: %s", err)
	}

	shovel, err := chs.Consume("page", "shovel", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume source: %s", err)
	}


	///////////////////////////////////////

	// Setup the destination of the store and forward
	destination, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("connection.open destination: %s", err)
	}
	defer destination.Close()

	chd, err := destination.Channel()
	if err != nil {
		log.Fatalf("channel.open destination: %s", err)
	}

	if err := chd.ExchangeDeclare("log2", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("exchange.declare destination: %s", err)
	}


	if _, err := chs.QueueDeclare("remote-tee", true, true, false, false, nil); err != nil {
		log.Fatalf("queue.declare destination: %s", err)
	}

	if err := chs.QueueBind("remote-tee", "#", "log2", false, nil); err != nil {
		log.Fatalf("queue.bind destination: %s", err)
	}

	messages, err := chs.Consume("remote-tee", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume destination: %s", err)
	}

	go func() {
		for msg := range messages {
			log.Printf("remote-tee: %s", string(msg.Body))
			msg.Ack(false)
		}
	}()


	// Buffer of 1 for our single outstanding publishing
	confirms := chd.NotifyPublish(make(chan amqp.Confirmation, 1))

	if err := chd.Confirm(false); err != nil {
		log.Fatalf("confirm.select destination: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()


	chs.Publish("log1", "alert", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte("Hello world"),
	})



	// Now pump the messages, one by one, a smarter implementation
	// would batch the deliveries and use multiple ack/nacks
	for {
		msg, ok := <-shovel
		if !ok {
			log.Fatalf("source channel closed, see the reconnect example for handling this")
		}

		log.Printf("page: %s\n",string(msg.Body))

		err = chd.PublishWithContext(ctx, "log2", msg.RoutingKey, false, false, amqp.Publishing{
		//err = chd.PublishWithContext(ctx, "log3", msg.RoutingKey, false, false, amqp.Publishing{
			// Copy all the properties
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			DeliveryMode:    msg.DeliveryMode,
			Priority:        msg.Priority,
			CorrelationId:   msg.CorrelationId,
			ReplyTo:         msg.ReplyTo,
			Expiration:      msg.Expiration,
			MessageId:       msg.MessageId,
			Timestamp:       msg.Timestamp,
			Type:            msg.Type,
			UserId:          msg.UserId,
			AppId:           msg.AppId,

			// Custom headers
			Headers: msg.Headers,

			// And the body
			Body: msg.Body,
		})

		if err != nil {
			if e := msg.Nack(false, false); e != nil {
				log.Printf("nack error: %+v", e)
			}
			log.Fatalf("basic.publish destination: %+v", msg)
		}

		// only ack the source delivery when the destination acks the publishing
		if confirmed := <-confirms; confirmed.Ack {
			log.Println("page:msg.Ack",confirmed.Ack)
			if e := msg.Ack(false); e != nil {
				log.Printf("ack error: %+v", e)
			}
		} else {
			log.Println("page:msg.Nack")
			if e := msg.Nack(false, false); e != nil {
				log.Printf("nack error: %+v", e)
			}
		}
	}

}