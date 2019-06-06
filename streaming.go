package main

import (
	"fmt"

	stan "github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "client-123"
	sc, _ := stan.Connect(clusterID, clientID)

	// Simple async subscriber.
	sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("received message: %s\n", string(m.Data))
	}, stan.DurableName("my-durable"))

	// Simple synchronous publisher.
	// This does not return until an ack has been received from
	// NATS Streaming.
	sc.Publish("foo", []byte("hello world"))

	sub.Unsubscribe()
	sc.Close()
}
