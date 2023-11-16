## RabbitMQ example of Event Driven Architecture with github.com/rabbitmq/amqp091-go \*

Setup:

```
go mod init eventdrivenrabbit

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
go get golang.org/x/sync/errgroup
```

Start:

```
go run .\cmd\consumer\main.go
go run .\cmd\producer\main.go
```
