package broker

import (
	"io"
	"log"
	"net"
	"sync"
)

type peer struct {
	broker     Broker
	connection net.Conn
	buffer     []byte
	lock       sync.Mutex
}

func (p *peer) Write(topic string, data []byte) {
	m := make([]byte, 2)
	m[0] = 'p'
	m[1] = topic[0]
	m = append(m, data[:4]...)
	if _, err := p.connection.Write(m); err != nil {
		log.Fatal(err)
	}
}

func (p *peer) Read() {
	defer p.close()

	data := make([]byte, 10)

	for {
		if n, err := p.connection.Read(data); err != nil {
			if err == io.EOF {
				continue
			}
			log.Fatal(err)
		} else {
			//log.Printf("Read %d bytes", n)
			p.handle(data[:n])
		}
	}
}

func (p *peer) handle(data []byte) {
	p.buffer = append(p.buffer, data...)
	//log.Println(string(p.buffer))

	p.lock.Lock()
	for len(p.buffer) >= 6 {
		switch p.buffer[0] {
		case 'p':
			//log.Printf("Push message '%s' in topic '%s'", p.buffer[2:6], p.buffer[1:2])
			p.broker.PeerWrite(string(p.buffer[1:2]), p.buffer[2:6])
			p.buffer = p.buffer[6:]
		default:
			break
		}
	}
	p.lock.Unlock()
}

func (p *peer) close() {
	if err := p.connection.Close(); err != nil {
		log.Println(err)
	}
}

func newPeer(broker Broker, connection net.Conn) peer {
	return peer{
		broker:     broker,
		connection: connection,
		buffer:     make([]byte, 0),
		lock:       sync.Mutex{},
	}
}
