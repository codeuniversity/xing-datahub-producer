package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/codeuniversity/xing-datahub-producer/handler"
	"github.com/codeuniversity/xing-datahub-producer/metrics"
	"github.com/codeuniversity/xing-datahub-protocol"
	"github.com/prometheus/client_golang/prometheus"

	"golang.org/x/crypto/acme/autocert"
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
	initPrometheus()

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
		ProtoMessage: &protocol.TargetUser{},
		Topic:        "target_users",
	}
	http.Handle("/target_users", targetUserHandler)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}

	//ssl handling

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/certs"), //Folder for storing certificates
	}
	server := &http.Server{
		Addr: ":3000",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
}

func initPrometheus() {
	prometheus.MustRegister(metrics.HTTPProcessed)

	http.Handle("/metrics", prometheus.Handler())
}
