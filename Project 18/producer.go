package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)
	

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	amqpServerURL  := os.Getenv("AMQP_SERVER_URL")
	if amqpServerURL  == "" {
		amqpServerURL  = "amqp://localhost"
	}

	// Create a new RabbitMQ connection.
    connectRabbitMQ, err := amqp.Dial(amqpServerURL )
    if err != nil {
        panic(err)
    }
    defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
    // instance over the connection we have already
    // established.
    channelRabbitMQ, err := connectRabbitMQ.Channel()
    if err != nil {
        panic(err)
    }
    defer channelRabbitMQ.Close()

    // With the instance and declare Queues that we can
    // publish and subscribe to.
    _, err = channelRabbitMQ.QueueDeclare(
        "QueueService1", // queue name
        true,            // durable
        false,           // auto delete
        false,           // exclusive
        false,           // no wait
        nil,             // arguments
    )
    if err != nil {
        panic(err)
    }

    // Create a new Fiber instance.
    app := fiber.New()

    // Add middleware.
    app.Use(
        logger.New(), // add simple logger
    )

    // Add route.
    app.Get("/send", func(c *fiber.Ctx) error {
        // Create a message to publish.
        message := amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(c.Query("msg")),
        }

        // Attempt to publish a message to the queue.
        if err := channelRabbitMQ.Publish(
            "",              // exchange
            "QueueService1", // queue name
            false,           // mandatory
            false,           // immediate
            message,         // message to publish
        ); err != nil {
            return err
        }

        return nil
    })

    // Start Fiber API server.
    log.Fatal(app.Listen(":3000"))
}
