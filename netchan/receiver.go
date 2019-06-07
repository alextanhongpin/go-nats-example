package main

import (
	"fmt"
	"log"

	nats "github.com/nats-io/go-nats"
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

	recvCh := make(chan *Person)
	ec.BindRecvChan("hello", recvCh)
	who := <-recvCh
	fmt.Println("got:", who)
}
