package main

import (
	"log"
	"time"
	"math/rand"
	"net"
	"encoding/binary"
	"flag"
	"encoding/base64"
	"network"
)

const (
	PING = 1 << iota
	PONG
)

type Queue string

// ... consumes
type Consumer interface {
	MessageWriter
}

type toLogConsumer struct {
	prefix string
}

func (t *toLogConsumer) Write(message Message) {
	log.Printf("%s%s", t.prefix, message.body)
}

// Wrapper for all the nodes in the cluster
type Broker interface {
	MessageWriter

	Subscribe(queue Queue, consumer Consumer)
	Join(node Node)
	Leave(node Node)
	MarkListen(queue Queue, id uint)
}

type broker struct {
	nodeMap  map[uint]Node  // Each node is allocated at a fixed index, 32 max
	queueMap map[Queue]uint // Each queue has a mask of known nodes that want these messages
}

func (b *broker) MarkListen(queue Queue, id uint) {
	b.queueMap[queue] = b.queueMap[queue] | id
	log.Printf("Queue %s has mask %b", queue, b.queueMap[queue])
}

func (b *broker) Join(node Node) {
	for index, n := range b.nodeMap {
		if n == nil {
			b.nodeMap[index] = node
			node.SetId(uint(index + 1))
			log.Printf("Node joined and was assigned id %d", index)
			return
		}
	}
	panic("Node pool is full")
}

func (b *broker) Leave(node Node) {
	id := node.GetId()
	delete(b.nodeMap, id)
}

func (b *broker) Subscribe(queue Queue, consumer Consumer) {
	values := make([]Node, 0)
	for _, n := range b.nodeMap {
		if n == nil {
			continue
		}
		values = append(values, n)
	}
	target := values[rand.Int()%len(values)]
	target.Subscribe(queue, consumer)
}

func NewBroker() Broker {
	const n = 32
	nodes := make(map[uint]Node, n)
	for i := 0; i < n; i++ {
		nodes[uint(i)] = nil
	}

	return &broker{nodes, make(map[Queue]uint, 0)}
}

func (b *broker) Write(message Message) {
	targets := b.queueMap[message.queue]
	for i := 0; i < 32; i++ {
		if (targets & (1 << uint(i))) != 0 {
			b.nodeMap[uint(i)].Write(message)
		}
	}
}

type NodeMessage struct {
	messageType uint8
	body        []byte
}

type NodeProxy interface {
	Send(message NodeMessage)
}

type physicalNode struct {
	inBuffer  []byte
	outBuffer []byte
}

func (p *physicalNode) Read() []byte {
	t := p.outBuffer
	p.outBuffer = make([]byte, 0)
	//log.Printf("Will send to socket %X", t)
	return t
}

func (p *physicalNode) Init(addr net.Addr) {
	log.Printf("Physical node connected from: %s", addr.String())
}

func (p *physicalNode) Write(bytes []byte) {
	p.inBuffer = append(p.inBuffer, bytes...)

	if binary.BigEndian.Uint32(p.inBuffer) == 0xCAFEBABE {
		size := binary.BigEndian.Uint32(p.inBuffer[4:8])
		bytes := p.inBuffer[8:8+size]
		p.Handle(NodeMessage{
			bytes[0],
			bytes[1:],
		})
	}
}

func (p *physicalNode) Closed() {
	log.Println("Physical node disconnected")
}

func (p *physicalNode) Send(message NodeMessage) {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(0xCAFEBABE))
	n := len(message.body) + 1
	b := make([]byte, 5)
	binary.BigEndian.PutUint32(b, uint32(n))
	b[4] = message.messageType
	b = append(b, message.body...)
	header = append(header, b...)

	p.outBuffer = append(p.outBuffer, header...)
}

func (p *physicalNode) Handle(message NodeMessage) {
	switch message.messageType {
	case PING:
		log.Println("I got pinged")
		p.Send(NodeMessage{PONG,make([]byte, 0)})
		break
	case PONG:
		log.Println("I got a PONG response")
	case 3:
		log.Printf("Execute: %s", message.body)

		break
	default:
		log.Println("Unknown messate type")
	}
}

// Proxy for a physical node
type Node interface {
	MessageWriter

	Subscribe(queue Queue, consumer Consumer)
	GetId() uint
	SetId(id uint)
}

type virtualNode struct {
	broker    Broker
	consumers []Consumer
	id        uint
}

func (n *virtualNode) SetId(id uint) {
	n.id = id
}

func (n *virtualNode) GetId() uint {
	return n.id
}

func (n *virtualNode) Subscribe(queue Queue, consumer Consumer) {
	n.consumers = append(n.consumers, consumer)
	log.Printf("Subscribed on node %d", n.id)
	n.broker.MarkListen(queue, 1<<(n.id-1))
}

func (n *virtualNode) Write(message Message) {
	for _, c := range n.consumers {
		c.Write(message)
	}
}

type Message struct {
	queue Queue
	body  []byte
}

type MessageWriter interface {
	Write(message Message)
}

func connectToPeer(address string) {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Println("Problem connecting")
	}

	log.Println("Connected to peer")
	node := &physicalNode{}
	go network.NewTcpLoop(conn, node)

	// Just send a ping
	node.Send(NodeMessage{
		PING,
		make([]byte, 0),
	})
}

func main() {

	// Setup thy rand
	rand.Seed(int64(time.Now().Nanosecond()))

	// Get flags
	listen := flag.String("listen", "0.0.0.0:2222", "")
	join := flag.String("join", "", "")

	flag.Parse()

	token := base64.StdEncoding.EncodeToString([]byte(*listen))

	log.Printf("Others can join with the token: %s", token)

	server := network.NewTcpServer()
	go server.Start(*listen, func() network.TcpClient {
		return &physicalNode{
			inBuffer:  make([]byte, 0),
			outBuffer: make([]byte, 0),
		}
	})

	if join != nil && len(*join) > 0 {
		log.Printf("We should join others at %s", *join)
		peer, err := base64.StdEncoding.DecodeString(*join)
		if err != nil {
			panic("Error decoding join token")

		}
		go connectToPeer(string(peer))
	}

	time.Sleep(time.Hour)
}
