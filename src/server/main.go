package main

import (
	"network"
	"net"
	"log"
	"bytes"
)

func wrapConnection(conn net.Conn) network.Protocol {
	proto := network.NewProtocol()

	go network.CreateLoop(conn, proto.BytesOut(), proto.BytesIn())
	go proto.Run()

	return proto
}

func connectionHandler(conn net.Conn) {
	proto := wrapConnection(conn)

	for {
		m := <- proto.MessagesIn()
		switch m.Kind {
		case network.PING:
			log.Println("We got a ping, better respond")
			proto.MessagesOut() <- network.Message{network.PONG, nil}
		case network.REVERSE:
			proto.MessagesOut() <- network.Message{network.REVERSE, bytes.ToUpper(m.Body)}
		default:
			log.Println("Unknown message kind")
		}
	}
}

func main() {
	tcpServer := network.NewTcpServer()
	tcpServer.Start("0.0.0.0:3333", connectionHandler)
}
