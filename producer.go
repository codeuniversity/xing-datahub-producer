package main

import (
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/codeuniversity/xing-datahub-producer/Protocol"
	"github.com/codeuniversity/xing-datahub-producer/handler"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	brokerList := []string{"localhost:9092"}
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		panic(err)
	}
	userHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &Protocol.User{},
		Topic:        "users",
	}
	http.Handle("/users", userHandler)

	connectionHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &Protocol.Connection{},
		Topic:        "connections",
	}
	http.Handle("/connections", connectionHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
