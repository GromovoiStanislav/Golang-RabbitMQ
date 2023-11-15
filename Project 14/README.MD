## RabbitMQ with github.com/rabbitmq/amqp091-go example\*

Setup:

```
go mod init go-rabbitmq-example

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

Start:

```
go run .\consumer\consumer.go log.WARN log.ERROR
go run .\consumer\consumer.go log.*


go run .\producer\producer.go log.WARN
go run .\producer\producer.go log.ERROR
go run .\producer\producer.go log.INFO
```
