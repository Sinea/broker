package p2p

import "encoding/binary"

func buildMessage(from, to PeerID, body []byte) []byte {
	data := make([]byte, 10+len(body))

	// Write header + meta
	binary.BigEndian.PutUint16(data[0:2], 0xCAFE)
	binary.BigEndian.PutUint16(data[2:4], uint16(from))
	binary.BigEndian.PutUint16(data[4:6], uint16(to))
	binary.BigEndian.PutUint32(data[6:10], uint32(len(body)))

	// Copy body
	copy(data[10:], body)

	return data
}
