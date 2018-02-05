package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// RequestHandler serializes json posts to protobuf messages and pushes them to the specified kafka topic
type RequestHandler struct {
	Producer     sarama.AsyncProducer
	Topic        string
	ProtoMessage proto.Message
}

func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := checkToken(r); err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	h.ProtoMessage.Reset()
	jsonpb.Unmarshal(r.Body, h.ProtoMessage)
	message, err := proto.Marshal(h.ProtoMessage)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	m := &sarama.ProducerMessage{
		Topic: h.Topic,
		Value: sarama.ByteEncoder(message),
	}

	h.Producer.Input() <- m
	w.WriteHeader(200)
}

func checkToken(r *http.Request) error {
	envToken := os.Getenv("token")
	if envToken == "" {
		return nil
	}
	token := r.Header.Get("access-token")
	if token != envToken {
		return errors.New("acess-token doesn't match, got: " + token)
	}
	return nil
}
