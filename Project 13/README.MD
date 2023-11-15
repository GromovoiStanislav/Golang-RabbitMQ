## RabbitMQ example of amqp.Confirmation and Exchange/QueueInspect with github.com/rabbitmq/amqp091-go

Setup:

```
go mod init main

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

Start:

```
go run producer.go
go run inspection.go
```
