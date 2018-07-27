package wire

import (
	"testing"
	"log"
)

func TestPacking(t *testing.T) {
	packed := Pack([]byte{1, 2, 3})
	log.Printf("%X", packed)

	unpacked := NewUnpacker()(packed)
	log.Printf("%X", unpacked)
}
