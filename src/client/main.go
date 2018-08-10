package main

import (
	"net"
	"github.com/sinea/network/client"
	"github.com/sinea/network/wire"
	"time"
	"github.com/sinea/network/io"
)

func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", "0.0.0.0:3333")
	if err != nil {
		panic("Shit")
	}

	return conn
}

func main() {

	conn := connectToServer()
	messageWriter := client.NewMessageWriter(wire.NewWriter(io.NewWriter(conn)))
	for {
		messageWriter.Write(client.Hello{
			Name: "bile",
			Age:  30,
		})
		time.Sleep(time.Second)
	}

	time.Sleep(time.Hour)
}
