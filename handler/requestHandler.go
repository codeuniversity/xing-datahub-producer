package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/codeuniversity/xing-datahub-producer/metrics"
	protocol "github.com/codeuniversity/xing-datahub-protocol"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// RequestHandler serializes json posts to protobuf messages and pushes them to the specified kafka topic
type RequestHandler struct {
	Producer   sarama.AsyncProducer
	Topic      string
	RawMessage protocol.RawMessage
}

func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			h.answerWith(w, http.StatusInternalServerError)
		}
	}()

	if err := checkToken(r); err != nil {
		fmt.Println(err)
		h.answerWith(w, http.StatusUnauthorized)
		return
	}

	h.RawMessage.Reset()
	if err := jsonpb.Unmarshal(r.Body, h.RawMessage); err != nil {
		h.answerWith(w, http.StatusBadRequest)
		return
	}

	message, err := proto.Marshal(h.RawMessage.Parse())
	if err != nil {
		h.answerWith(w, http.StatusInternalServerError)
		return
	}
	m := &sarama.ProducerMessage{
		Topic: h.Topic,
		Value: sarama.ByteEncoder(message),
	}

	h.Producer.Input() <- m
	h.answerWith(w, http.StatusOK)
}

func (h *RequestHandler) answerWith(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	metrics.HTTPProcessed.WithLabelValues(strconv.Itoa(code), h.Topic).Inc()
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
