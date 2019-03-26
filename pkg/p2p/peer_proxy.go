package p2p

type peerProxy struct {
	peerID PeerID
}

// Send data via a route
func (p *peerProxy) Send(data []byte) {
	// Route the message through the connected peers
}
