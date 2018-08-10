package main

import (
	"net"
	"github.com/sinea/network/client"
	"time"
	"log"
	"github.com/sinea/network/p2p"
)

func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", "0.0.0.0:3333")
	if err != nil {
		panic("Shit")
	}

	return conn
}

func wratata(writer client.MessageWriter, i interface{}) {
	for {
		writer.Write(i)
		time.Sleep(time.Second)
	}
}

func main() {
	log.Println("Starting client")
	//conn := connectToServer()
	//messageWriter := client.NewMessageWriter(wire.NewWriter(io.NewWriter(conn)))
	//go wratata(messageWriter, client.Hello{"Me", 3})
	//go wratata(messageWriter, client.Goodbye{":)"})
	//
	//time.Sleep(time.Hour)

	c := p2p.NewClient()

	if err := c.Connect("0.0.0.0:3333"); err != nil {
		panic(err)
	}

	reply := <-c.Request(p2p.NewAskRequest(":)"), p2p.RequestOptions{Timeout: 10})

	if r, ok := reply.(p2p.AskReply); ok {
		log.Printf("Received reply %s", r.Result)
	}

	c.Message(p2p.PeerAuth{"123", "3333"})
}
