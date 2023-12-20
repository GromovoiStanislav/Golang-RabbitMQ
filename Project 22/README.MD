## RabbitMQ Creating multiple instances of the consumer with github.com/rabbitmq/amqp091-go

#### Setup:

```
go mod init main

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

#### Start:

```
go run .\consumer\
go run .\producer\
```
