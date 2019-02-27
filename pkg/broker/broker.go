package broker

import (
	"fmt"
	"hash/crc32"
	"sync"
)

const MaxShards = 64

type MessageWriter interface {
	Write(topic string, data []byte)
}

type MessageReader interface {
	Read(topic string) <-chan []byte
}

type Broker interface {
	MessageWriter
	MessageReader
}

type broker struct {
	nodeCount        uint32 // Total number of nodes in cluster
	topicReplication uint32 // Known number of replicas for each topic
	shardMask        uint64 // Shard bit mask for the current node

	shards map[uint64]shard
	topics map[string][]chan []byte
	lock   *sync.Mutex
}

func (b *broker) Read(topic string) <-chan []byte {
	b.lock.Lock()
	out := make(chan []byte)
	if _, ok := b.topics[topic]; !ok {
		b.topics[topic] = make([]chan []byte, 0)
	}
	b.topics[topic] = append(b.topics[topic], out)
	b.lock.Unlock()

	return out
}

func (b *broker) Write(topic string, data []byte) {
	// The data should be written to this shard
	shardID := b.getShardID(topic)

	fmt.Printf("Sending %s to shard %d\n", data, shardID)
	if out, ok := b.topics[topic]; ok {
		for _, o := range out {
			o <- data
		}
	}
}

func (b *broker) getShardID(topic string) uint64 {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(topic), table)
	// The data should be written to this shard
	return 1 << ((checksum % b.topicReplication) % MaxShards)
}

func New(shardCount uint32) Broker {
	return &broker{
		lock:             &sync.Mutex{},
		topicReplication: shardCount,
		topics:           map[string][]chan []byte{},
	}
}
