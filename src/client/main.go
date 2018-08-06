package main

import (
	"net"
	"time"
	"github.com/sinea/network/io"
	"github.com/sinea/network/wire"
	"github.com/sinea/network/client"
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
	ioWriter := io.NewWriter(conn)
	wireWriter := wire.NewWriter(ioWriter)
	messageWriter := client.NewMessageWriter(wireWriter)
	for {
		messageWriter.Write(client.NewMessage(5, []byte{0xCA, 0xFE}))
		time.Sleep(time.Second)
	}

	time.Sleep(time.Hour)
}
