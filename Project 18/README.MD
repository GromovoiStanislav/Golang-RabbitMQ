## RabbitMQ example with github.com/rabbitmq/amqp091-go

This is a demo application using Fiber and RabbitMQ

#### Setup:

```
go mod init main

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
go get github.com/gofiber/fiber/v2
```

#### Start:

```
go run consumer.go
```

```
go run producer.go
```
