package main

import (
	"log"
	"time"
	"math/rand"
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
	MarkListen(queue Queue, id int)
}

type broker struct {
	nodeMap  map[int]Node  // Each node is allocated at a fixed index, 32 max
	queueMap map[Queue]int // Each queue has a mask of known nodes that want these messages
}

func (b *broker) MarkListen(queue Queue, id int) {
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
	panic("implement me")
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
	nodes := make(map[int]Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = nil
	}

	return &broker{nodes, make(map[Queue]int, 0)}
}

func (b *broker) Write(message Message) {
	targets := b.queueMap[message.queue]
	for i := 0; i < 32; i++ {
		if (targets & (1 << uint(i))) != 0 {
			b.nodeMap[i].Write(message)
		}
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

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	broker := NewBroker()

	q1 := Queue("a.b.c")
	q2 := Queue("x.y.z")

	m1 := Message{q1,[]byte("hello")}
	m2 := Message{q2,[]byte("bye")}

	for i := 0; i < 32; i++ {
		consumers := make([]Consumer, 0, 0)
		broker.Join(&virtualNode{
			broker,
			consumers,
			0,
		})
	}

	broker.Subscribe(q1, &toLogConsumer{"a: "})
	broker.Subscribe(q1, &toLogConsumer{"b: "})
	broker.Subscribe(q2, &toLogConsumer{"c: "})

	broker.Write(m1)
	broker.Write(m2)

}
