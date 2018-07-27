package main

import (
	"network"
	"net"
	"log"
)

func connectionHandler(conn net.Conn) {
	log.Println("New client connected")
	client := network.NewClient(conn)

	for {
		select {
		case message := <-client.Messages():

			switch message.Kind {
			case network.PING:
				client.Reply(network.Message{Kind: network.PONG}, message)
				break
			}

			break
		}
	}
}

func main() {
	network.Listen("0.0.0.0:3333", connectionHandler)
}
