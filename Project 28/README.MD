## RabbitMQ stream example with github.com/rabbitmq/amqp091-go

This repository is an example of how streaming works in RabbitMQ in Go.

#### Setup:

```
go mod init rabbit-example

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

#### Start:

```
go run ./producer
go run ./consumer
```
