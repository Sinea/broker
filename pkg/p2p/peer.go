package p2p

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type peer struct {
	messageID  uint32
	peerID     uint16
	fromID     uint16
	connection net.Conn
	messages   chan<- Message
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
	if n, err := p.connection.Write(packedMessage); err == nil {
		log.Printf("Wrote %d bytes", n)
	} else {
		log.Fatal(err)
	}
}

// Route route this data through this node
func (p *peer) Route(to uint16, data []byte) {
	log.Printf("Routing to %d message: %s", to, data)
}

func (p *peer) handle(message []byte) {
	to := binary.BigEndian.Uint16(message[4:6])
	if to != p.peerID {
		p.Route(to, message)
	} else {
		from := binary.BigEndian.Uint16(message[2:4])
		p.messages <- Message{from, message[10:]}
	}
}

func newPeer(connection net.Conn, messages chan<- Message, id uint16) *peer {
	return &peer{
		connection: connection,
		messages:   messages,
		fromID:     id,
	}
}
