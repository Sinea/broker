package broker

import (
	"log"
	"net"
	"sync"
)

type MessageWriter interface {
	Write(topic string, data []byte)
}

type MessageReader interface {
	Read(topic string) <-chan []byte
}

type Broker interface {
	MessageWriter
	MessageReader

	Start(address string) error
	Join(address string) error
	PeerWrite(topic string, data []byte)
}

type broker struct {
	uncheckedPeers []*peer
	topics         map[string]chan []byte
	lock           *sync.Mutex
}

func (b *broker) Join(address string) error {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	b.handle(connection)

	return nil
}

func (b *broker) Start(address string) error {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	defer closeListener(listener)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		b.handle(connection)
	}
}

func (b *broker) Read(topic string) <-chan []byte {
	b.lock.Lock()
	if _, ok := b.topics[topic]; !ok {
		b.topics[topic] = make(chan []byte)
	}
	b.lock.Unlock()

	return b.topics[topic]
}

func (b *broker) Write(topic string, data []byte) {
	// Send to local listeners
	if out, ok := b.topics[topic]; ok {
		out <- data
	}

	// Send to all peers
	for _, p := range b.uncheckedPeers {
		p.Write(topic, data)
	}
}

func (b *broker) PeerWrite(topic string, data []byte) {
	if out, ok := b.topics[topic]; ok {
		out <- data
	}
}

func closeListener(listener net.Listener) {
	if err := listener.Close(); err != nil {
		log.Println(err)
	}
}

func (b *broker) handle(connection net.Conn) {
	p := newPeer(b, connection)
	b.uncheckedPeers = append(b.uncheckedPeers, &p)
	go p.Read()
}

func New(shardCount uint32) Broker {
	return &broker{
		lock:           &sync.Mutex{},
		topics:         make(map[string]chan []byte),
		uncheckedPeers: make([]*peer, 0),
	}
}
