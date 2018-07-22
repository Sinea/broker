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

type TcpServer interface {
	Start(addr string, factory SocketHandler)
}

type tcpServer struct {
}

func (s *tcpServer) Start(addr string, socketHandler SocketHandler) {
	log.Printf("Listening on %s", addr)
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		log.Printf("Error listening on %s", addr)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting new client: %s", err.Error())
		} else {
			go socketHandler(conn)
		}
	}
}

func NewTcpServer() TcpServer {
	return &tcpServer{}
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

func NewProtocol() Protocol {
	return &proto{
		make(chan Message),
		make(chan Message),
		make(chan []byte),
		make(chan []byte),
		make([]byte, 0),
	}
}

type Protocol interface {
	MessagesIn() <-chan Message
	MessagesOut() chan<- Message

	BytesIn() chan<- []byte
	BytesOut() <-chan []byte

	Run()
}

type Message struct {
	Kind uint8
	Body []byte
}

type proto struct {
	msgIn    chan Message
	msgOut   chan Message
	bytesIn  chan []byte
	bytesOut chan []byte
	buffer   []byte
}

func (p *proto) BytesIn() chan<- []byte {
	return p.bytesIn
}

func (p *proto) BytesOut() <-chan []byte {
	return p.bytesOut
}

func (p *proto) MessagesIn() <-chan Message {
	return p.msgIn
}

func (p *proto) MessagesOut() chan<- Message {
	return p.msgOut
}

func (p *proto) Run() {
	for {
		select {
		case m := <-p.msgOut:
			p.pack(m)
		case b := <-p.bytesIn:
			p.buffer = append(p.buffer, b...)
			p.unpack()
		default:
			time.Sleep(42 * time.Millisecond)
		}
	}
}

// Write a message to the bytesOut channel
func (p *proto) pack(message Message) {
	buffer := []byte{0xCE, uint8(len(message.Body) + 1), message.Kind}
	buffer = append(buffer, message.Body...)
	log.Printf("Write bytes: %X", buffer)
	p.bytesOut <- buffer
}

// Read from bytes and send to msgIn
func (p *proto) unpack() {
	if p.buffer[0] == 0xCE { // We have a valid message header
		size := uint8(p.buffer[1])
		if uint8(len(p.buffer)-2) < size { // Not all data was received
			return
		}
		p.msgIn <- Message{p.buffer[2], p.buffer[3:]}
		p.buffer = p.buffer[2+size:]
	}
}
