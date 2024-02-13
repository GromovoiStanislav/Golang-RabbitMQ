package main

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpURL := os.Getenv("CLOUDAMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}
	
	conn, err := amqp.DialConfig(amqpURL, amqp.Config{
		Properties: amqp.NewConnectionProperties(),
	})
	if err != nil {
		panic(err)
	}

	var channel *amqp.Channel
	channel, err = conn.Channel()
	if err != nil {
		panic(err)
	}

	///start dlx
	var restingQueue amqp.Queue
	restingQueue, err = channel.QueueDeclare("my-dl-queue", true, false, false, false, map[string]interface{}{
		"x-dead-letter-exchange":    "my-exchange",
		"x-dead-letter-routing-key": "my-routing-key",
		"x-max-length":              50000,
		"x-overflow":                "reject-publish",
		"x-message-ttl":             1000,
	})
	if err != nil {
		panic(err)
	}

	err = channel.ExchangeDeclare("my-dlx", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(restingQueue.Name, "my-routing-key", "my-dlx", false, nil)
	if err != nil {
		panic(err)
	}
	///end dlx


	err = channel.ExchangeDeclare("my-exchange", "direct", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	var queue amqp.Queue
	queue, err = channel.QueueDeclare("my-queue", true, false, false, false, map[string]interface{}{
		"x-dead-letter-exchange": "my-dlx",
	})
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(queue.Name, "my-routing-key", "my-exchange", false, nil)
	if err != nil {
		panic(err)
	}

	var deliveries <-chan amqp.Delivery
	deliveries, err = channel.Consume(queue.Name, "my-consumer", false, false, false, false, nil)


	for msg := range deliveries {

		go func(msg amqp.Delivery) {
		
			if err := processMsg(msg); err != nil {
				if msg.Headers["x-death"] == nil{
					_ = msg.Nack(false, false)
				}

				if msg.Headers["x-death"] != nil {
					for _, death := range msg.Headers["x-death"].([]interface{}) {
						deathMap := death.(amqp.Table)
						if deathMap["reason"] == "expired" {
							count, ok := deathMap["count"].(int64)
							if ok && count < 5 {
								_ = msg.Nack(false, false)
							} else {
								fmt.Printf("maximum retries has been exceeded: %v\n", msg.MessageId)
								msg.Ack(true)
							}
							break
						}
					}
				}
				
			} else {
				msg.Ack(false)
			}

		}(msg)
	}

}

func processMsg(msg amqp.Delivery) error {
	fmt.Println(string(msg.Body))
	
	return errors.New("Ошибка")
	//return nil
}