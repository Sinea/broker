package p2p

import "errors"

type peerProxy struct {
	remote PeerID
	local  PeerID
	m      *mesh
}

// this should never happen
func (p *peerProxy) write(data []byte) error {
	// Nothing happens here
	return errors.New("this should be unreachable")
}

// Send data via a route
func (p *peerProxy) Send(data []byte) {
	t := buildMessage(p.local, p.remote, 0, data)
	p.m.sendToPeer(p.remote, t)
}

// create a new peer proxy
func newPeerProxy(local, remote PeerID, m *mesh) Peer {
	return &peerProxy{remote, local, m}
}
