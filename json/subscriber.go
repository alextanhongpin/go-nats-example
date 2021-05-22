package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	var wg sync.WaitGroup
	wg.Add(2)

	c.Subscribe("foo", func(s string) {
		defer wg.Done()
		fmt.Println("received from foo:", s)
	})
	type Person struct {
		Name string
		Age  int
	}
	c.Subscribe("bar", func(p *Person) {
		defer wg.Done()
		fmt.Println("received from bar:", p)
	})
	wg.Wait()
}
