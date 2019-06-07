package main

import (
	"fmt"
	"log"

	nats "github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.ChanSubscribe("hello", ch)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	msg := <-ch
	fmt.Println("received msg", string(msg.Data))
}
