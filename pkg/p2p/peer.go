package p2p

import (
	"encoding/binary"
	"fmt"
	"net"
)

type peerProxy struct {
	peerID uint16
}

func (*peerProxy) Send(data []byte) {
	// TODO : Route this through a connected peer
}

type peer struct {
	messageID  uint32
	peerID     uint16
	connection net.Conn
}

func (p *peer) Send(data []byte) {
	// TODO : Write this to the underlying socket
	_, _ = p.connection.Write(data)
}

func (p *peer) sendMessage(message []byte) {
	// TODO : Send the message using the routing table
}

func (p *peer) handle(message []byte) {
	target := binary.BigEndian.Uint16(message[4:6])
	if target != p.peerID { // Message is not for me
		p.sendMessage(message)
	} else { // Message is for me
		// TODO : Check if the received message is not a duplicate then handle the message
		fmt.Println(string(message[10:]))
	}
}

func (p *peer) read() {
	// TODO : Read data from the underlying socket
	data := make([]byte, 100)
	n, _ := p.connection.Read(data)
	p.handle(data[:n])
}

func buildMessage(from, to uint16, id uint32, body []byte) []byte {
	data := make([]byte, 10+len(body))
	// Write header + meta
	binary.BigEndian.PutUint16(data[0:2], 0xCAFE)
	binary.BigEndian.PutUint16(data[2:4], from)
	binary.BigEndian.PutUint16(data[4:6], to)
	binary.BigEndian.PutUint32(data[6:10], id)
	// Copy body
	copy(data[10:], body)

	return data
}
