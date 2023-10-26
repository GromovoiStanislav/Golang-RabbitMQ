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

## Example #4: Routing:

```
go run receive_logs_direct.go info warn
go run emit_log_direct.go warn "a warning"
```

Direct exchange:

<image src="img/direct-exchange.png" alt="preview">

Multiple bindings:

<image src="img/direct-exchange-multiple.png" alt="preview">

Putting it all together:

<image src="img/four.png" alt="preview">

## Example #5: Topics:

```
go run receive_logs_topic.go "#"                    - To receive all the logs
go run receive_logs_topic.go "kern.*"               - To receive all logs from the facility "kern"
go run receive_logs_topic.go "*.critical"           - if you want to hear only about "critical" logs
go run receive_logs_topic.go "kern.*" "*.critical"  - You can create multiple bindings

go run emit_log_topic.go kern.critical "A critical kernel error"
```

<image src="img/topic-exchange.png" alt="preview">

```
* (star) can substitute for exactly one word.
# (hash) can substitute for zero or more words.

When a queue is bound with "#" (hash) binding key - it will receive all the messages, regardless of the routing key - like in 'fanout' exchange.

When special characters "*" (star) and "#" (hash) aren't used in bindings, the topic exchange will behave just like a 'direct' one.
```
