package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	args := os.Args[1:]
	pushOrPull := args[0]

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to connect to nats: %s", err)
	}

	// Create JetStream Context.
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatalf("failed to create jetstream context: %s", err)
	}

	var sub *nats.Subscription
	switch pushOrPull {
	case "push":
		// Simple Async Ephemeral Consumer.
		js.Subscribe("ORDERS.*", func(m *nats.Msg) {
			fmt.Printf("received a JetStream message: %v\n", m)
		})

	case "push_durable":
		// Simple Sync Durable Consumer (optional SubOpts at the end)
		var err error
		sub, err = js.SubscribeSync("ORDERS.scratch", nats.Durable("SOME_RANDOM_NAME"), nats.MaxDeliver(3))
		//sub, err = js.SubscribeSync("ORDERS.scratch", nats.MaxDeliver(3))
		if err != nil {
			log.Fatalf("failed to subscribe sync: %s", err)
		}
		for {
			timeout := 5 * time.Second
			m, err := sub.NextMsg(timeout)
			if err != nil {
				log.Fatalf("failed to read next msg: %s", err)
			}
			log.Printf("received message: %+v", m)
			if err := m.Ack(); err != nil {
				log.Fatalf("failed to ack message: %v", err)
			}
		}
	case "pull":
		// Simple Pull Consumer.
		sub, err = js.PullSubscribe("ORDERS.scratch", "NEW")
		if err != nil {
			log.Fatalf("failed to subscribe pull: %s", err)
		}

		for {
			msgs, err := sub.Fetch(10)
			if err != nil {
				log.Fatalf("failed to fetch message: %s", err)
			}
			log.Printf("received messages: %v\n", msgs)
			for _, msg := range msgs {
				if err := msg.Ack(); err != nil {
					log.Fatalf("failed to ack message: %v", err)
				}
			}
			if len(msgs) == 0 {
				break
			}
		}
		log.Println("done fetching all messages")
	case "chan":
		ch := make(chan *nats.Msg)
		// Simple Pull Consumer.
		sub, err = js.ChanSubscribe("ORDERS.scratch", ch, nats.Durable("WHY"))
		if err != nil {
			log.Fatalf("failed to subscribe pull: %s", err)
		}
		for {
			select {
			case msg := <-ch:
				log.Printf("received msg: %v\n", msg)
				if err := msg.Ack(); err != nil {
					log.Fatalf("failed to ack message: %v", err)
				}
			}
		}

	default:
		log.Fatalf(`invalid command %q: must be "push" or "pull"`, pushOrPull)
	}

	// Unsubscribe.
	sub.Unsubscribe()

	// Drain.
	sub.Drain()
}
