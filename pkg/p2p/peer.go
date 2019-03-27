package p2p

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	ReadBufferSize = 1024

	// Message flags
	isSystemMessage byte = 1 << iota
	isHandshake
)

type peer struct {
	remote     PeerID
	local      PeerID
	connection net.Conn
	messages   chan<- Message
	m          *mesh
	buffer     []byte
}

// Read data from the socket and try to handle or route the data
func (p *peer) Read() {
	buffer := make([]byte, ReadBufferSize)
	for {
		if n, err := p.connection.Read(buffer); err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		} else {
			if err := p.handle(buffer[:n]); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Send data via socket
func (p *peer) Send(data []byte) error {
	packedMessage := buildMessage(p.local, p.remote, 0, data)
	if err := p.write(packedMessage); err != nil {
		return err
	}

	return nil
}

// raw write to socket
func (p *peer) write(message []byte) error {
	if n, err := p.connection.Write(message); err != nil {
		return fmt.Errorf("error writing to socket: %s", err.Error())
	} else if n < len(message) {
		return errors.New("message not written")
	}

	return nil
}

// handle a packed message
func (p *peer) handle(message []byte) error {
	p.buffer = append(p.buffer, message...)
	// We don't have sufficient data
	if len(p.buffer) < 8 {
		return nil
	}

	// All messages start with 0xCE byte
	if p.buffer[0] != 0xCE {
		return fmt.Errorf("invalid message header: %x", p.buffer[0])
	}

	// Now let's check if the message is complete
	messageSize := binary.BigEndian.Uint32(p.buffer[4:8])
	if uint32(len(p.buffer)) < 8+messageSize {
		return nil // We don't have the complete message
	}

	// Check if the message is for the local node
	dst := PeerID(p.buffer[3])
	if dst != p.remote {
		// Route the message to appropriate peer
		p.m.sendToPeer(dst, p.buffer[:8+messageSize])
	} else {
		// Message is for the local node
		src := PeerID(p.buffer[2])
		flags := p.buffer[1]
		body := p.buffer[8:]
		p.buffer = p.buffer[8+messageSize:]
		// Make something with the data
		if (flags & isSystemMessage) != 0 {
			log.Printf("Handle system messsage: %s", string(body))
			p.handleSystemMessage(src, flags, body)
		} else {
			p.messages <- Message{src, body}
		}
	}

	return nil
}

func (p *peer) handleSystemMessage(src PeerID, flags byte, body []byte) {
	if (flags & isHandshake) != 0 {
		log.Println("This is a handshake")

		m := IdExchangeMessage{}
		if err := json.Unmarshal(body, &m); err != nil {
			log.Println(err)
		}
		log.Printf("Hello %d", m.Id)
		p.local = m.Id
		p.m.peers[m.Id] = p
		p.m.nodes[m.Id] = p
	}
}

// create a new peer
func newPeer(connection net.Conn, messages chan<- Message, m *mesh) *peer {
	return &peer{
		connection: connection,
		messages:   messages,
		buffer:     make([]byte, 0),
		m:          m,
	}
}
