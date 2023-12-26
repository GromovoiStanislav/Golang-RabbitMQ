## RabbitMQ AlternateExchange example with github.com/rabbitmq/amqp091-go

#### Setup:

```
go mod init main

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

#### Start:

```
go run producer-1.go
go run producer-2.go
```
