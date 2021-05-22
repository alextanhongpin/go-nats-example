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
		log.Fatal(err)
	}
	defer nc.Close()

	msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("got message", msg)
}
