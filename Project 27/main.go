package main

import (
	"fmt"

	"rabbit-example/internal/rabbitmq"
)

type App struct {
	Rmq *rabbitmq.RabbitMQ
}

func Run() error {
	fmt.Println("Go RabbitMQ simple example")

	rmq := rabbitmq.NewRabbitMQService()
	app := App{
		Rmq: rmq,
	}

	err := app.Rmq.Connect()
	if err != nil {
		return err
	}
	defer app.Rmq.Conn.Close()

	err = app.Rmq.Publish("Hi")
	if err != nil {
		return err
	}

	forever := make(chan bool)
	app.Rmq.Consume()
	<-forever

	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Println("Error Setting Up our application")
		fmt.Println(err)
	}
}