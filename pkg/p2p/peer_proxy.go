package p2p

type peerProxy struct {
	peerID uint16
}

// Send data via a route
func (*peerProxy) Send(data []byte) {
}
