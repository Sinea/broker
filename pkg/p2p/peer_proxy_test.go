package p2p

import "testing"

func TestPeerProxy_Send(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should have panicked")
		}
	}()
	proxy := newPeerProxy(0, 1, nil)
	proxy.write([]byte{})
}
