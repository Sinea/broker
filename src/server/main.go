package main

import (
	"network"
	"net"
)

func connectionHandler(conn net.Conn) {
	client := network.NewClient(conn)

	for {
		select {
		case message := <-client.Messages():

			switch message.Kind {
			case network.PING:
				client.Send(network.Message{Kind:network.PONG})
				break
			}

			break
		}
	}
}

func main() {
	network.Listen("0.0.0.0:3333", connectionHandler)
}
