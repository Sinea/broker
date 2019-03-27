package p2p

import (
	"fmt"
	"testing"
)

func TestPeer_Broadcast(t *testing.T) {
	p := &peer{peerID: 5}
	p.Send([]byte{1, 2, 3, 4})
	p.Send([]byte{1, 2, 3, 4})
	p.Send([]byte{1, 2, 3, 4})
	p.Send([]byte{1, 2, 3, 4})
}

func TestPeerProxy_Sendx(t *testing.T) {
	var mesh Mesh
	mesh.Listen("0.0.0.0:1111")

	mesh.Broadcast([]byte{1, 2, 3, 4})

	if p, err := mesh.Peer(3); err != nil {
		p.Send([]byte{1, 2, 3, 4})
	}

	for message := range mesh.Read() {
		fmt.Printf("Received %s from %d", string(message.Data), message.From)
	}
}

func TestPeer_Route(t *testing.T) {
	fmt.Printf("%d", buildMessage(1, 2, []byte("sinea")))
}
