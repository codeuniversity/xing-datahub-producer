package main

import (
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/codeuniversity/xing-datahub-producer/Protocol"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func makeMessageFor(m proto.Message, topic string, r *http.Request) (*sarama.ProducerMessage, error) {
	jsonpb.Unmarshal(r.Body, m)
	message, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}, nil
}

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	brokerList := []string{"localhost:9092"}
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		user := &Protocol.User{}
		kafkaMessage, err := makeMessageFor(user, "users", r)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		producer.SendMessage(kafkaMessage)
		w.WriteHeader(200)
	})

	http.HandleFunc("/connections", func(w http.ResponseWriter, r *http.Request) {
		c := &Protocol.Connection{}
		kafkaMessage, err := makeMessageFor(c, "connections", r)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		producer.SendMessage(kafkaMessage)
		w.WriteHeader(200)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
