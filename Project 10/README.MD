## Official example RabbitMQ with github.com/rabbitmq/amqp091-go

This exports a Client object that wraps this library. It automatically reconnects when the connection fails, and blocks all pushes until the connection succeeds. It also confirms every outgoing message, so none are lost. It doesn't automatically ack each message, but leaves that to the parent process, since it is usage-dependent.

Try running this in one terminal, and `rabbitmq-server` in another. Stop & restart RabbitMQ to see how the queue reacts.

setup:

```
go mod init main

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

start:

```
go run consumer.go
go run producer.go
```
