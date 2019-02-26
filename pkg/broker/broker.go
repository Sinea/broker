package broker

import (
	"fmt"
	"hash/crc32"
)

const MaxShards = 64

type MessageWriter interface {
	Write(topic string, data []byte)
}

type Broker interface {
	MessageWriter
}

type broker struct {
	nodeCount uint32 // Total number of nodes in cluster
	numShards uint32 // Known number of shards
	shardMask uint64 // Shard bit mask for the current node

	shards map[uint64]shard
}

func (b *broker) Write(topic string, data []byte) {
	// The data should be written to this shard
	shardID := b.getShardID(topic)

	fmt.Printf("Sending %s to shard %d\n", data, shardID)
}

func (b *broker) getShardID(topic string) uint64 {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(topic), table)
	// The data should be written to this shard
	return 1 << ((checksum % b.numShards) % MaxShards)
}

func New(shardCount uint32) Broker {
	return &broker{
		numShards: shardCount,
	}
}
