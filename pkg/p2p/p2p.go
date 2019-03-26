package p2p

type PeerID uint8

type Message struct {
	From PeerID
	Data []byte
}

type Peer interface {
	Send(data []byte)
}

type Mesh interface {
	Listen(address string) error
	Join(address string) error
	Broadcast(data []byte)
	Peer(ID PeerID) (Peer, error)
	Read() <-chan Message
}
