package event

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestEmitterCreateSuccess(t *testing.T) {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Errorf("Could not establish a connection to AMQP server: %v", err)
	}
	defer conn.Close()
	_, err = NewEventEmitter(conn)
	if err != nil {
		t.Errorf("Error creating Event Emitter: %v", err)
	}
}

func TestEmitterPushSuccess(t *testing.T) {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Errorf("Could not establish a connection to AMQP server: %v", err)
	}

	defer conn.Close()
	emitter, err := NewEventEmitter(conn)
	if err != nil {
		t.Errorf("Error creating Event Emitter: %v", err)
	}

	err = emitter.Push("Hello World!", "INFO")
	if err != nil {
		t.Errorf("Could not push to queue successfully: %v", err)
	}
}

func TestConsumerCreatSuccess(t *testing.T) {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Errorf("Could not establish a connection to AMQP server: %v", err)
	}

	defer conn.Close()
}