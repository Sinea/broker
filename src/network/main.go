package network

import (
	"log"
	"net"
)

type TcpClientFactory func() TcpClient

type TcpClient interface {
	// Handle the incoming bytes
	Write(bytes []byte)

	// Socket is already closed. Perform some cleanup!?
	Closed()
}

type TcpServer interface {
	Start(addr string, factory TcpClientFactory)
}

type tcpServer struct {
}

func (s *tcpServer) Start(addr string, factory TcpClientFactory) {
	log.Printf("Listening on %s", addr)
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		log.Printf("Error listening on %s", addr)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting new client")
		}
		// Spawn a new goroutine to handle the client (for now)
		go func(conn net.Conn) {
			client := factory()
			for {
				d := make([]byte, 1024)
				c, err := conn.Read(d)

				if err != nil {
					log.Println("Error reading from socket")
					client.Closed()
					break
				}
				client.Write(d[:c])
			}

		}(conn)
	}
}

func NewTcpServer() TcpServer {
	return &tcpServer{}
}
