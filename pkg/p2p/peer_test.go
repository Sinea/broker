package p2p

import "testing"

func TestPeer_handle_low_data(t *testing.T) {
	peer := newPeer(nil, make(chan Message), 0)
	err := peer.handle([]byte{0xCE})
	if err != nil {
		t.Fatal("No error expected")
	}
}

func TestPeer_write_wrong_header(t *testing.T) {
	peer := newPeer(nil, make(chan Message), 0)
	err := peer.handle([]byte{0xCF, 0, 0, 0, 0, 0, 0, 0, 0})
	if err == nil {
		t.Fatal("No error expected")
	}
}

func TestPeer_write_not_enaugh_data(t *testing.T) {
	peer := newPeer(nil, make(chan Message), 0)
	err := peer.handle([]byte{0xCE, 0, 0, 0, 0, 0, 0, 100})
	if err != nil {
		t.Fatal("No error expected")
	}
}
