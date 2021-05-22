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
	ack, err := js.Publish("ORDERS.scratch", []byte("hello"))
	if err != nil {
		log.Fatalf("failed to publish: %s", err)
	}
	log.Println("ack", ack)

	// Simple Async Stream Publisher.
	for i := 0; i < 500; i++ {
		js.PublishAsync("ORDERS.scratch", []byte("hello"))
	}
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}

	// Simple Async Ephemeral Consumer.
	js.Subscribe("ORDERS.*", func(m *nats.Msg) {
		fmt.Printf("received a JetStream message: %s\n", string(m.Data))
	})

	// Simple Sync Durable Consumer (optional SubOpts at the end)
	sub, err := js.SubscribeSync("ORDERS.*", nats.Durable("MONITOR"), nats.MaxDeliver(3))
	if err != nil {
		log.Fatalf("failed to subscribe sync: %s", err)
	}
	timeout := 5 * time.Second
	m, err := sub.NextMsg(timeout)
	if err != nil {
		log.Fatalf("failed to read next msg: %s", err)
	}
	log.Printf("received message: %v\n", m)

	// Simple Pull Consumer.
	sub, err = js.PullSubscribe("ORDERS.*", "MONITOR")
	if err != nil {
		log.Fatalf("failed to subscribe pull: %s", err)
	}

	msgs, err := sub.Fetch(10)
	if err != nil {
		log.Fatalf("failed to fetch message: %s", err)
	}

	log.Printf("received messages: %v\n", msgs)

	// Unsubscribe.
	sub.Unsubscribe()

	// Drain.
	sub.Drain()
}
