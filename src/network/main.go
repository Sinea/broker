package network

import (
	"log"
	"net"
	"time"
	"errors"
)

const (
	PING    = 1 << iota
	PONG
	REVERSE
)

type ByteBufferHandler func(buffer []byte)
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

var tag uint8 = 0

type Message struct {
	Kind uint8
	tag  uint8
}

type Client interface {
	Request(message Message) <-chan Message
	Reply(reply Message, to Message)
	Send(message Message)
	Messages() <-chan Message
}

type mockClient struct {
	requests map[uint8]chan Message
	messages chan Message
	bytesOut chan []byte
}

func (m *mockClient) Reply(reply Message, to Message) {
	reply.tag = to.tag
	m.Send(reply)
}

func (m *mockClient) Send(message Message) {
	m.bytesOut <- pack(message)
}

func (m *mockClient) Request(message Message) <-chan Message {
	c := make(chan Message)
	message.tag = tag
	tag = (tag + 1) % 0xFF
	m.requests[message.tag] = c
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
		make(chan []byte),
	}

	// Read
	go func(c net.Conn, msgIn chan Message) {
		buffer := make([]byte, 0)
		for {
			b := make([]byte, 1024)
			c.SetReadDeadline(time.Now().Add(time.Millisecond))
			n, err := c.Read(b)
			if err != nil {
				if terr, ok := err.(net.Error); ok {
					if !terr.Timeout() {
						log.Printf("Not a timeout error: %s", terr)
						break
					}
				} else {
					log.Printf("Error reading from socket: %s", err)
					break
				}
			} else {
				buffer = append(buffer, b[:n]...)
				m := Message{}
				if e := unpack(buffer, &m); e != nil {
					log.Println("Error unpacking")
					break
				}
				buffer = buffer[3:]

				if ch, ok := client.requests[m.tag]; ok {
					ch <- m
					close(ch)
					delete(client.requests, m.tag)
				} else {
					client.messages <- m
				}
			}
		}
	}(conn, client.messages)

	// Write
	go func(c net.Conn, bts <-chan []byte) {
		xyz := make([]byte, 0)
		for {
			select {
			case b := <-bts:
				xyz = append(xyz, b...)
				break
			default:
				if len(xyz) == 0 {
					time.Sleep(time.Millisecond)
					continue
				}
				c.SetWriteDeadline(time.Now().Add(time.Millisecond))
				_, err := c.Write(xyz)
				if err != nil {
					if terr, ok := err.(net.Error); ok {
						if !terr.Timeout() {
							log.Printf("Not a timeout error: %s", terr)
							break
						}
					} else {
						log.Printf("Error reading from socket: %s", err)
						break
					}
				} else {
					xyz = make([]byte, 0)
				}
			}
		}
	}(conn, client.bytesOut)

	return client
}

func unpack(buffer []byte, message *Message) error {
	if buffer[0] == 0xEE {
		message.Kind = buffer[1]
		message.tag = buffer[2]
		return nil
	}

	return errors.New(":(")
}

func pack(message Message) []byte {
	return []byte{0xEE, message.Kind, message.tag}
}
