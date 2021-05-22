package main

import (
	"log"

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
	// Publish string.
	c.Publish("foo", "hello world")
	type Person struct {
		Name string
		Age  int
	}
	// Publish struct.
	c.Publish("bar", &Person{"john", 20})
}
