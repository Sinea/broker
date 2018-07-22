package main

import (
	"net"
	"log"
	"network"
)

type EchoClient interface {
	Ping()
}

type echoClient struct {
	in  chan []byte
	out chan []byte
}

func (c *echoClient) Ping() {
	c.in <- []byte("SEND NUDES")
}

func (c *echoClient) Run() {
	for {
		m := <-c.out
		log.Printf("Received: %s", m)
	}
}

func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", "0.0.0.0:3333")
	if err != nil {
		log.Printf("Error connecting to server: %s", err.Error())
		panic("Shit")
	} else {
		return conn
	}
}

func wrapConnection(conn net.Conn) network.Protocol {
	proto := network.NewProtocol()

	go network.CreateLoop(conn, proto.BytesOut(), proto.BytesIn())
	go proto.Run()

	return proto
}

func main() {
	conn := connectToServer()
	defer conn.Close()
	proto := wrapConnection(conn)

	proto.MessagesOut() <- network.Message{network.PING, nil}
	reply := <-proto.MessagesIn()
	log.Printf("Received reply: %d", reply.Kind)

	proto.MessagesOut() <- network.Message{network.REVERSE, []byte("bile")}
	reply = <-proto.MessagesIn()
	log.Printf("Received reply: %s", reply.Body)
}
