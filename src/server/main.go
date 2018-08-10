package main

import (
	"log"
	"net"
	"github.com/sinea/network/client"
	"github.com/sinea/network/wire"
	"github.com/sinea/network/io"
)

func Handle(m interface{}) {
	switch m.(type) {
	case *client.Hello:
		h := m.(*client.Hello)
		log.Printf("Hello from %s", h.Name)
		break
	case *client.Goodbye:
		g := m.(*client.Goodbye)
		log.Printf("Goodbye %d", g.From)
		break
	default:
		panic("unknown type")
	}
}

func main() {

	listener, _ := net.Listen("tcp", "0.0.0.0:3333")

	conn, _ := listener.Accept()

	messages := client.NewReader()
	w := wire.NewReader(messages)
	io.NewReader(conn, w)

	for {
		Handle(<-messages.Messages())
	}
}
