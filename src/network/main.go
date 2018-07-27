package network

import (
	"log"
	"net"
	"time"
)

const (
	PING    = 1 << iota
	PONG
	REVERSE
)

type SocketHandler func(conn net.Conn)

func Listen(address string, handler SocketHandler) error {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Printf("Error listening on %s", address)
		return err
	}

	log.Printf("Listening on %s", address)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting new client: %s", err.Error())
			return err
		} else {
			go handler(conn)
		}
	}
}

func CreateLoop(conn net.Conn, in <-chan []byte, out chan<- []byte) {
	go func(conn net.Conn, out chan<- []byte) {
		for {
			// Read from socket and send to channel
			buffer := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond))
			n, err := conn.Read(buffer)
			if err != nil {
				if terr, ok := err.(net.Error); ok {
					if !terr.Timeout() {
						log.Printf("Not a timeout error: %s", terr)
					}
				} else {
					log.Printf("Error reading from socket: %s", err)
					break
				}
			} else {
				out <- buffer[:n]
			}
		}
	}(conn, out)

	for {
		// Read from channel and write to socket
		buffer := <-in
		conn.SetWriteDeadline(time.Now().Add(time.Millisecond))
		n, err := conn.Write(buffer)
		if err != nil {
			if terr, ok := err.(net.Error); ok {
				if !terr.Timeout() {
					log.Printf("Not a timeout error: %s", terr)
				}
			} else {
				log.Printf("Error reading from socket: %s", err)
				break
			}
		} else {
			log.Printf("Wrote %d bytes: %X", n, buffer)
		}
	}
}

var tag uint8 = 0

type Message struct {
	Kind uint8
	tag uint8
}

type Client interface {
	Request(message Message) <-chan Message
	Send(message Message)
	Messages() <-chan Message
}

type mockClient struct {
	requests map[uint8]chan Message
	messages chan Message
	bytesOut []byte
}

func (m *mockClient) Send(message Message) {
	log.Printf("Send: %d", message.Kind)
	buffer := []byte{0xEE, message.Kind, message.tag}
	m.bytesOut = append(m.bytesOut, buffer...)
	log.Printf("On wire: %X", m.bytesOut)
}

func (m *mockClient) Request(message Message) <-chan Message {
	c := make(chan Message)
	message.tag = tag
	m.requests[message.tag] = c

	// Replace with async
	go func(t uint8) {
		time.Sleep(time.Second - time.Millisecond*10)
		m.requests[t] <- Message{PONG, t}
		close(m.requests[t])
		delete(m.requests, t)
	}(message.tag)
	tag = (tag + 1) % 0xFF

	m.Send(message)

	return c
}

func (m *mockClient) Messages() <-chan Message {
	return m.messages
}

func NewClient(conn net.Conn) Client {

	client := &mockClient{
		make(map[uint8]chan Message),
		make(chan Message),
		make([]byte, 0),
	}

	// Read
	go func(c net.Conn, msgIn chan Message) {
		buffer := make([]byte, 0)
		for {
			b := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond))
			n, err := conn.Read(buffer)
			if err != nil {
				if terr, ok := err.(net.Error); ok {
					if !terr.Timeout() {
						log.Printf("Not a timeout error: %s", terr)
					}
				} else {
					log.Printf("Error reading from socket: %s", err)
					break
				}
			} else {
				buffer = append(buffer, b...)
				decode(buffer, msgIn)
			}
		}
	}(conn, client.messages)

	// Write

	return client
}

func decode(buffer []byte, out chan Message) {
	if buffer[0] == 0xEE {
		m := Message{buffer[1], buffer[2]}
		buffer = buffer[2:]
		out <- m
	}
}