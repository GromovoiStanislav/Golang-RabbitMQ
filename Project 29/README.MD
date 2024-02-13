## Retrying Mechanism for RabbitMQ and Dead Letter Queue example with github.com/rabbitmq/amqp091-go

#### Setup:

```
go mod init rabbit-example

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

#### Start:

```
go run consumer.go
go run producer.go
```
