package p2p

import (
	"encoding/binary"
	"log"
)

func buildMessage(from, to PeerID, flags byte, body []byte) []byte {
	data := make([]byte, 8+len(body))

	// Write header + meta
	data[0] = 0xCE
	data[1] = flags
	data[2] = byte(from)
	data[3] = byte(to)

	// Write data length
	binary.BigEndian.PutUint32(data[4:8], uint32(len(body)))

	// Copy body
	copy(data[8:], body)

	return data
}

func buildSystemMessage(from, to PeerID, flags, kind byte, v interface{}) []byte {
	body := bytes(v)
	log.Println(string(body))
	bodySize := uint32(len(body) + 1)
	data := make([]byte, 9+len(body))

	// Write header + meta
	data[0] = 0xCE
	data[1] = flags
	data[2] = byte(from)
	data[3] = byte(to)
	data[8] = kind

	// Write data length
	binary.BigEndian.PutUint32(data[4:8], uint32(bodySize))

	// Copy body
	copy(data[9:], body)

	return data
}
