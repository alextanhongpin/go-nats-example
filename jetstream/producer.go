package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to connect to nats: %s", err)
	}

	// Create JetStream Context.
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatalf("failed to create jetstream context: %s", err)
	}

	// Simple Stream Publisher.
	ack, err := js.Publish("ORDERS.scratch", []byte("hello world"))
	if err != nil {
		log.Fatalf("failed to publish: %s", err)
	}
	log.Println("ack", ack)

	// Simple Async Stream Publisher.
	for i := 0; i < 500; i++ {
		js.PublishAsync("ORDERS.scratch", []byte(fmt.Sprintf("hello, %s", time.Now())))
	}
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}
}
