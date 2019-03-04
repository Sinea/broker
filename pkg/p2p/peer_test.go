package p2p

import (
	"testing"
)

func TestPeer_Broadcast(t *testing.T) {
	p := &peer{peerID: 5}
	p.Send(3, []byte{1, 2, 3, 4})
	p.Send(3, []byte{1, 2, 3, 4})
	p.Send(3, []byte{1, 2, 3, 4})
	p.Send(3, []byte{1, 2, 3, 4})
}
