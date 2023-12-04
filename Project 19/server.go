package main

import (
	"net/http"

	"rabbitmq-example/route"
)

func main() {
	http.HandleFunc("/", route.Routeing())

	http.ListenAndServe(":8080", nil)
}