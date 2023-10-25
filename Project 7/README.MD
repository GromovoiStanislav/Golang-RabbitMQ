# Official examples using RabbitMQ with github.com/rabbitmq/amqp091-go

## Setup:

```
go mod init main

go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

## Example #1: "Hello World!":

```
go run send.go
go run receive.go
```

<image src="img/one.png" alt="preview">

## Example #2: Work Queues:

```
go run new_task.go hello world.....
go run worker.go
go run worker.go
```

<image src="img/two.png" alt="preview">

## Example #3: Publish/Subscribe:

```
go run receive_logs.go
go run receive_logs.go
go run emit_log.go hello world
```

<image src="img/exchanges.png" alt="preview">
