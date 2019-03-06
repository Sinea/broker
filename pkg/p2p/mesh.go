package p2p

import (
	"log"
	"net"
)

type Peer interface {
	Send(data []byte)
}

type Mesh interface {
	Listen(address string)
	Join(address string)
	Broadcast(data []byte)
	Peer(ID uint16) Peer
}

type mesh struct {
	peers map[uint16]Peer // Only connected peers
	nodes map[uint16]Peer // All nodes
}

func (m *mesh) Listen(address string) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal(err)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// TODO : Handle new connection
		log.Println("New connection from: " + connection.RemoteAddr().String())
	}
}

func (m *mesh) Join(address string) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to: %s", connection.RemoteAddr().String())
}

func (m *mesh) Broadcast(data []byte) {
	for _, node := range m.nodes {
		node.Send(data)
	}
}

func (m *mesh) Peer(ID uint16) Peer {
	return m.peers[ID]
}

func New() Mesh {
	return &mesh{}
}
