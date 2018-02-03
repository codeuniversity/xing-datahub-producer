package main

import (
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/codeuniversity/xing-datahub-producer/handler"
	"github.com/codeuniversity/xing-datahub-protocol"
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
		ProtoMessage: &protocol.User{},
		Topic:        "users",
	}
	http.Handle("/users", userHandler)
	connectionHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &protocol.Connection{},
		Topic:        "connections",
	}
	http.Handle("/connections", connectionHandler)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
