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

	sub, err := nc.SubscribeSync("hello")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	m, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(m.Data))
	fmt.Printf("%#v", m)
}
