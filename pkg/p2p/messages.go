package p2p

type PeerListMessage struct {
	Peers []PeerID `json:"peers"`
}

type IdExchangeMessage struct {
	Id PeerID `json:"id"`
}
