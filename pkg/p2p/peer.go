package p2p

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

const isSystemMessage byte = 1 << iota

type peer struct {
	peerID     PeerID
	fromID     PeerID
	connection net.Conn
	messages   chan<- Message
	m          *mesh
	buffer     []byte
}

// Read data from the socket and try to handle or route the data
func (p *peer) Read() {
	buffer := make([]byte, 1024)
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
func (p *peer) Send(data []byte) {
	packedMessage := buildMessage(p.fromID, p.peerID, data)
	if err := p.write(packedMessage); err != nil {
		log.Fatal(err.Error())
	}
}

// raw write to socket
func (p *peer) write(message []byte) error {
	if n, err := p.connection.Write(message); err != nil {
		return err
	} else {
		log.Printf("Wrote %d bytes", n)
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
	if dst != p.peerID {
		// Route the message to appropriate peer
		p.m.sendToPeer(dst, p.buffer[:8+messageSize])
	} else {
		// Message is for the local node
		src := PeerID(p.buffer[2])
		flags := p.buffer[1]
		body := p.buffer[10:]
		p.buffer = p.buffer[8+messageSize:]
		// Make something with the data
		if (flags & isSystemMessage) != 0 {
			log.Printf("Handle system messsage: %s", string(body))
		} else {
			p.messages <- Message{src, body}
		}
	}

	return nil
}

// create a new peer
func newPeer(connection net.Conn, messages chan<- Message, id PeerID) *peer {
	return &peer{
		connection: connection,
		messages:   messages,
		fromID:     id,
		buffer:     make([]byte, 0),
	}
}
