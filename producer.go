package main

import (
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/codeuniversity/xing-datahub-producer/handler"
	"github.com/codeuniversity/xing-datahub-protocol"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	brokerList := []string{"localhost:9092"}
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		panic(err)
	}
	userHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &protocol.User{},
		Topic:        "users",
	}
	http.Handle("/users", userHandler)

	itemHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &protocol.Item{},
		Topic:        "items",
	}
	http.Handle("/items", itemHandler)

	interactionHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &protocol.Interaction{},
		Topic:        "interactions",
	}
	http.Handle("/interactions", interactionHandler)

	targetItemHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &protocol.TargetItem{},
		Topic:        "target_items",
	}
	http.Handle("/target_items", targetItemHandler)

	targetUserHandler := handler.RequestHandler{
		Producer:     producer,
		ProtoMessage: &protocol.Interaction{},
		Topic:        "target_users",
	}
	http.Handle("/target_users", targetUserHandler)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
