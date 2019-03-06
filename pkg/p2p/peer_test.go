package p2p

import (
	"testing"
)

func TestPeer_Broadcast(t *testing.T) {
	p := &peer{peerID: 5}
	p.Send([]byte{1, 2, 3, 4})
	p.Send([]byte{1, 2, 3, 4})
	p.Send([]byte{1, 2, 3, 4})
	p.Send([]byte{1, 2, 3, 4})
}

func TestPeerProxy_Send(t *testing.T) {
	var mesh Mesh
	mesh.Listen("0.0.0.0:1111")

	mesh.Broadcast([]byte{1, 2, 3, 4})

	mesh.Peer(3).Send([]byte{1, 2, 3, 4})
}
