package p2p

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type peer struct {
	peerID     PeerID
	fromID     PeerID
	connection net.Conn
	messages   chan<- Message
	m          *mesh
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
			p.handle(buffer[:n])
		}
	}
}

// Send data via socket
func (p *peer) Send(data []byte) {
	packedMessage := buildMessage(p.fromID, p.peerID, data)
	p.write(packedMessage)
}

// raw write to socket
func (p *peer) write(message []byte) {
	if n, err := p.connection.Write(message); err == nil {
		log.Printf("Wrote %d bytes", n)
	} else {
		log.Fatal(err)
	}
}

// handle a packed message
func (p *peer) handle(message []byte) {
	to := PeerID(binary.BigEndian.Uint16(message[4:6]))
	if to != p.peerID {
		p.m.sendToPeer(to, message)
	} else {
		from := PeerID(binary.BigEndian.Uint16(message[2:4]))
		p.messages <- Message{from, message[10:]}
	}
}

// create a new peer
func newPeer(connection net.Conn, messages chan<- Message, id PeerID) *peer {
	return &peer{
		connection: connection,
		messages:   messages,
		fromID:     id,
	}
}
