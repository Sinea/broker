package p2p

type peerProxy struct {
	remote PeerID
	local  PeerID
	m      *mesh
}

// this should never happen
func (p *peerProxy) write(data []byte) {
	// Nothing happens here
	panic(`this should be unreachable`)
}

// Send data via a route
func (p *peerProxy) Send(data []byte) {
	t := buildMessage(p.local, p.remote, data)
	p.m.sendToPeer(p.remote, t)
}

// create a new peer proxy
func newPeerProxy(local, remote PeerID, m *mesh) Peer {
	return &peerProxy{remote, local, m}
}
