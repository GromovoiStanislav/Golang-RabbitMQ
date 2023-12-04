package rabbitmqconnect

import (
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"

	errors "rabbitmq-example/exceptions"
)

func ConnectMQ() (*amqp.Connection, *amqp.Channel) {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		errors.FailOnError(err, "Error loading .env file")
	}

	amqpUrl := os.Getenv("CLOUDAMQP_URL")
	if amqpUrl == "" {
		amqpUrl = "amqp://rabbitmq:rabbitmq@127.0.0.1:5672"
	}


	conn, err := amqp.Dial(amqpUrl)
	errors.FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	ch, err := conn.Channel()
	errors.FailOnError(err, "Failed to open a channel")
	// defer ch.Close()
	return conn, ch
}

func CloseMQ(conn *amqp.Connection, channel *amqp.Channel) {
	defer conn.Close()    //rabbit mq close
	defer channel.Close() //rabbit mq channel close
}