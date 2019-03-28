package p2p

import (
	"fmt"
	"log"
	"net"
)

type mesh struct {
	// Local ID
	ID        PeerID
	isRunning bool

	peers    map[PeerID]Peer // Only connected peers
	nodes    map[PeerID]Peer // All nodes
	messages chan Message

	routingTable map[PeerID]PeerID
}

// Read channel with the messages aimed at this node
func (m *mesh) Read() <-chan Message {
	return m.messages
}

// Listen on this address for peer connections
func (m *mesh) Listen(address string) (err error) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	log.Printf("Node %d is listening on %s", m.ID, address)

	defer func() {
		if err := listener.Close(); err != nil {
			log.Print(err.Error())
		}
	}()

	m.isRunning = true

	for m.isRunning {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("New connection from: " + connection.RemoteAddr().String())
		m.handleConnection(connection)
	}

	return
}

// Join by connecting to the provided address
func (m *mesh) Join(address string) (err error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	log.Printf("Connected to: %s", connection.RemoteAddr().String())
	m.handleConnection(connection)

	return
}

// Broadcast send the data to all nodes
func (m *mesh) Broadcast(data []byte) {
	for _, node := range m.nodes {
		if err := node.Send(data); err != nil {
			log.Println(err)
		}
	}
}

// Peer return a peer by it's id
func (m *mesh) Peer(ID PeerID) (Peer, error) {
	if peer, ok := m.nodes[ID]; ok {
		return peer, nil
	}

	return nil, fmt.Errorf("unknown peer with id %d", ID)
}

// route an already packed message to a remote peer
func (m *mesh) sendToPeer(remote PeerID, packedMessage []byte) {
	fmt.Printf("Send to %d data: %d", remote, packedMessage)
	routeID := m.routingTable[remote]
	if peer, err := m.Peer(routeID); err != nil {
		if err := peer.write(packedMessage); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

// handle a new connection
func (m *mesh) handleConnection(connection net.Conn) {
	peer := newPeer(connection, m.messages, m)
	go peer.Read()
	peer.initializeHandshake(m.ID)
}

func New(id PeerID) Mesh {
	return &mesh{
		ID:           id,
		peers:        map[PeerID]Peer{},
		nodes:        map[PeerID]Peer{},
		messages:     make(chan Message),
		routingTable: map[PeerID]PeerID{},
	}
}
