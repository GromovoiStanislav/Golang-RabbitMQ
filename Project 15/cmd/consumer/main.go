package main

import (
	"context"
	"log"
	"os"

	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"

	"eventdrivenrabbit/internal"
)

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost:5671"
	}

	conn, err := internal.ConnectRabbitMQ(url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
	publishConn, err := internal.ConnectRabbitMQ(url)
	defer publishConn.Close()

	mqClient, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	publishClient, err := internal.NewRabbitMQClient(publishConn)
	if err != nil {
		panic(err)
	}

	// Create Unnamed Queue which will generate a random name, set AutoDelete to True
	queue, err := mqClient.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	// CreateExchange
	if err := mqClient.CreateExchange("customer_events", "direct", true, true); err != nil {
		panic(err)
	}


	// Create binding between the customer_events exchange and the new Random Queue
	// Can skip Binding key since fanout will skip that rule
	if err := mqClient.CreateBinding(queue.Name, "customers.created.se", "customer_events"); err != nil {
		panic(err)
	}

	messageBus, err := mqClient.Consume(queue.Name, "email-service", false)
	if err != nil {
		panic(err)
	}

	// blocking is used to block forever
	var blocking chan struct{}

	// Set a timeout for 15 secs
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Create an Errgroup to manage concurrecy
	g, ctx := errgroup.WithContext(ctx)
	// Set amount of concurrent tasks
	g.SetLimit(10)

	// Apply Qos to limit amount of messages to consume
	if err := mqClient.ApplyQos(10, 0, true); err != nil {
		panic(err)
	}

	go func() {
		for message := range messageBus {
			// Spawn a worker
			msg := message
			g.Go(func() error {
				// Multiple means that we acknowledge a batch of messages, leave false for now
				if err := msg.Ack(false); err != nil {
					log.Printf("Acknowledged message failed: Retry ? Handle manually %s\n", msg.MessageId)
					return err
				}

				log.Printf("Acknowledged message, replying to %s\n", msg.ReplyTo)

				time.Sleep(5 * time.Second)

				// Use the msg.ReplyTo to send the message to the proper Queue
				if err := publishClient.Send(ctx, "customer_callbacks", msg.ReplyTo, amqp.Publishing{
					ContentType:   "text/plain",      // The payload we send is plaintext, could be JSON or others..
					DeliveryMode:  amqp.Transient, // This tells rabbitMQ to drop messages if restarted
					Body:          []byte("RPC Complete"),
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					panic(err)
				}
				return nil
			})
		}
	}()

	log.Println("Consuming, to close the program press CTRL+C")
	// This will block forever
	<-blocking

}