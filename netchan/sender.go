package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	sendCh := make(chan *Person)
	ec.BindSendChan("hello", sendCh)
	sendCh <- &Person{"john", 20}
}
