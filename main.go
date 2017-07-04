package main

import (
	"fmt"
	"time"

	"github.com/nats-io/go-nats"
)

type person struct {
	Name string
	Age  int
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer c.Close()
	if err != nil {
		panic(err)
	}

	// Simple async subscriber
	c.Subscribe("foo", func(p *person) {
		fmt.Printf("Received a message: name %s with age %d \n", p.Name, p.Age)
	})

	me := person{Name: "john.doe", Age: 1}
	// Simple publisher
	c.Publish("foo", me)

	time.Sleep(2 * time.Second)
	sub, err := c.Subscribe("foo", nil)
	sub.Unsubscribe()
	time.Sleep(2 * time.Second)
}
