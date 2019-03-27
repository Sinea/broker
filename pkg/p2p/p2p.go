package p2p

type PeerID uint8

type Message struct {
	From PeerID
	Data []byte
}

type Peer interface {
	writer
	Send(data []byte) error
}

type Mesh interface {
	Listen(address string) error
	Join(address string) error
	Broadcast(data []byte)
	Peer(ID PeerID) (Peer, error)
	Read() <-chan Message
}

type writer interface {
	write(data []byte) error
}
