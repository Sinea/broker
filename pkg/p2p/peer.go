package p2p

import (
	"encoding/binary"
	"fmt"
	"net"
)

var id uint32 = 0

type peer struct {
	messageID uint32
	peerID    uint16

	peers []peer
	nodes []peer
}

func (p *peer) Broadcast(data []byte) {
	// TODO : Send messages to all nodes through the connected peers
}

func (p *peer) Send(to uint16, body []byte) {
	message := buildMessage(p.peerID, to, p.messageID, body)
	p.messageID++
	fmt.Printf("%d\n", message)
	p.sendMessage(message)
}

func (p *peer) sendMessage(message []byte) {
	// TODO : Send the message using the routing table
}

func (p *peer) handle(message []byte) {
	if binary.BigEndian.Uint16(message[4:6]) != p.peerID { // Message is not for me
		p.sendMessage(message)
	} else { // Message is for me
		// TODO : Check if the received message is not a duplicate then handle the message
		fmt.Println(string(message[10:]))
	}
}

func Read(conn net.Conn) error {
	return nil
}

func buildMessage(from, to uint16, id uint32, body []byte) []byte {
	data := make([]byte, 10+len(body))
	// Write header + meta
	binary.BigEndian.PutUint16(data[0:2], 0xCAFE)
	binary.BigEndian.PutUint16(data[2:4], from)
	binary.BigEndian.PutUint16(data[4:6], to)
	binary.BigEndian.PutUint32(data[6:10], id)
	id++
	// Copy body
	copy(data[10:], body)

	return data
}
