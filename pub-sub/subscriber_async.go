package main

import (
	"fmt"
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
	nc.Subscribe("hello", func(m *nats.Msg) {
		defer wg.Done()
		fmt.Printf("receive msg: %s\n", string(m.Data))
	})
	wg.Wait()
}
