package main

import (
	"log"
	"github.com/sinea/network/client"
	"github.com/sinea/network/p2p"
)

func Handle(m interface{}) {
	switch m.(type) {
	case *client.Hello:
		h := m.(*client.Hello)
		log.Printf("Hello from %s", h.Name)
		break
	case *client.Goodbye:
		g := m.(*client.Goodbye)
		log.Printf("Goodbye %s", g.From)
		break
	default:
		panic("unknown type")
	}
}

func main() {
	log.Println("Starting server")
	//listener, _ := net.Listen("tcp", "0.0.0.0:3333")
	//
	//conn, _ := listener.Accept()
	//
	//messages := client.NewReader()
	//w := wire.NewReader(messages)
	//io.NewReader(conn, w)
	//
	//for {
	//	Handle(<-messages.Messages())
	//}

	c := p2p.NewClient()

	m := <-c.Messages()

	switch m.(type) {
	case p2p.AskRequest:
		m.(p2p.AskRequest).Reply(p2p.NewAskReply("it's ok"))
	}
}
