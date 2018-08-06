package main

import (
	"log"
	"net"
	"github.com/sinea/network/client"
	"github.com/sinea/network/wire"
	"github.com/sinea/network/io"
)

func main() {

	listener, _ := net.Listen("tcp", "0.0.0.0:3333")

	conn, _ := listener.Accept()

	messages := client.NewReader()
	w := wire.NewReader(messages)
	io.NewReader(conn, w)

	for {
		select {
		case m := <-messages.Messages():
			log.Printf("Received %d : %X", m.Kind(), m.Body())
			break
		}
	}
}
