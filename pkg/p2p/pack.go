package p2p

import "encoding/binary"

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
