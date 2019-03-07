package p2p

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	From uint16
	Data []byte
}

type Peer interface {
	PeerReader
	Send(data []byte)
}

type PeerReader interface {
	Read()
}

type MessageRouter interface {
	Route(to uint16, data []byte)
}

type Mesh interface {
	Listen(address string)
	Join(address string) error
	Broadcast(data []byte)
	Peer(ID uint16) (Peer, error)
	Read() <-chan Message
}

type mesh struct {
	// Local ID
	ID uint16

	peers    map[uint16]PeerReader // Only connected peers
	nodes    map[uint16]Peer       // All nodes
	messages chan Message
}

// Read channel with the messages aimed at this node
func (m *mesh) Read() <-chan Message {
	return m.messages
}

// Listen on this address for peer connections
func (m *mesh) Listen(address string) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Node %d is listening on %s", m.ID, address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("New connection from: " + connection.RemoteAddr().String())
		go m.handleConnection(connection)
	}
}

// Join by connecting to the provided address
func (m *mesh) Join(address string) (err error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	log.Printf("Connected to: %s", connection.RemoteAddr().String())
	go m.handleConnection(connection)

	return
}

// Broadcast send the data to all nodes
func (m *mesh) Broadcast(data []byte) {
	for _, node := range m.nodes {
		node.Send(data)
	}
}

// Peer return a peer by it's id
func (m *mesh) Peer(ID uint16) (Peer, error) {
	if peer, ok := m.nodes[ID]; ok {
		return peer, nil
	}

	return nil, fmt.Errorf("unknown peer with id %d", ID)
}

func New(id uint16) Mesh {
	return &mesh{
		ID:       id,
		peers:    map[uint16]PeerReader{},
		nodes:    map[uint16]Peer{},
		messages: make(chan Message),
	}
}
