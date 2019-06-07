package main

import (
	"log"
	"sync"

	nats "github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	nc.Subscribe("help", func(m *nats.Msg) {
		defer wg.Done()
		nc.Publish(m.Reply, []byte("I can help"))
	})
	wg.Wait()
}
