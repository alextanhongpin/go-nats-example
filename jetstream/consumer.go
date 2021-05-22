package main

import (
	"errors"
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
	defer nc.Drain()

	// Create JetStream Context.
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatalf("failed to create jetstream context: %s", err)
	}
	strInfo, err := js.StreamInfo("ORDERS")
	if err != nil {
		log.Fatalf("failed to get stream info: %v", err)
	}
	log.Printf("Stream already exists: %v\n", strInfo)

	conInfo, err := js.ConsumerInfo("ORDERS", "NEW")
	if err != nil {
		log.Fatalf("failed to get consumer info: %v", err)
	}
	log.Printf("Consumer already exists: %v\n", conInfo)

	// Create a Consumer
	//consumerInfo, err := js.AddConsumer("ORDERS", &nats.ConsumerConfig{
	//Durable:   "NEW",
	//AckPolicy: nats.AckExplicitPolicy,
	//})
	//if err != nil {
	//log.Fatalf("failed to create consumer: %v", err)
	//}
	//log.Printf("consumer info: %+v\n", consumerInfo)

	var sub *nats.Subscription
	switch pushOrPull {
	case "push":
		// Simple Async Ephemeral Consumer.
		js.Subscribe("ORDERS.scratch", func(m *nats.Msg) {
			fmt.Printf("received a JetStream message: %v\n", m)
		})

	case "push_durable":
		// Simple Sync Durable Consumer (optional SubOpts at the end)
		var err error
		sub, err = js.SubscribeSync("ORDERS.scratch", nats.Durable("NEW"), nats.MaxDeliver(3))
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
		// https://gist.github.com/wallyqs/b01ba613341170b4442acbffcaea0a81
		// Simple Pull Consumer.
		sub, err = js.PullSubscribe("ORDERS.scratch", "NEW")
		if err != nil {
			log.Fatalf("failed to subscribe pull: %s", err)
		}

		for {
			msgs, err := sub.Fetch(10)
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					log.Println("context exceeded when fetching")
					break
				}
				log.Printf("failed to fetch message: %s, sleeping for 1s", err)
				time.Sleep(1 * time.Second)
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
		sub, err = js.ChanSubscribe("ORDERS.scratch", ch, nats.Durable("NEW"))
		if err != nil {
			log.Fatalf("failed to subscribe pull: %s", err)
		}

		// Set limits of 1000 messages or 5MB, whichever comes first
		sub.SetPendingLimits(1000, 5*1024*1024)
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

	// Unsubscribe. (NOTE: This will remove the consumer from the stream).
	//sub.Unsubscribe()

	// Drain.
	sub.Drain()
}
