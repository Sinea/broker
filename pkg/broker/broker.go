package broker

import (
	"fmt"
	"hash/crc32"
)

type Writer interface {
	Write(topic string, data []byte)
}

type Broker interface {
	Writer
}

type broker struct {
	numShards uint32
}

func (b *broker) Write(topic string, data []byte) {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(topic), table)
	shard := checksum % b.numShards

	fmt.Printf("Sending %s to shard %d\n", data, shard)
}

func New() Broker {
	return &broker{
		numShards: 3,
	}
}
