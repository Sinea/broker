package p2p

import (
	"net"
)

// TODO : This should become a join negotiation
func (m *mesh) handleConnection(connection net.Conn) {
	//// Send my ID
	//log.Println("Sending my id")
	//b := make([]byte, 2)
	//binary.BigEndian.PutUint16(b, m.ID)
	//if n, err := connection.Write(b); err != nil {
	//	log.Fatal(err)
	//} else {
	//	log.Printf("Wrote %d bytes %d", n, b)
	//}
	//
	//// Read remote ID
	//log.Println("Waiting for remote")
	//data := make([]byte, 2)
	//if n, err := connection.Read(data); err != nil {
	//	log.Fatal(err)
	//} else if n < 2 {
	//	log.Fatal(fmt.Errorf("incomplete handshake byte count: %d", n))
	//} else {
	//	log.Printf("Received %d", data)
	//}
	//
	//// Add connected peer
	//peerID := binary.BigEndian.Uint16(data)
	//log.Printf("Peer with id %d joined", peerID)
	//m.nodes[peerID] = newPeer(connection, m.messages, m.ID)
	//go m.nodes[peerID].Read()
}
