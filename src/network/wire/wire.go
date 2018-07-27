package wire

import (
	"encoding/binary"
)

const (
	HEADER = 0xCA
	FOOTER = 0xBA
)

type wire struct {
	buffer []byte
}

func (w *wire) Unpack(chunk []byte) []byte {
	w.buffer = append(w.buffer, chunk...)
	l := len(w.buffer)

	if l < 5 { // Too few bytes
		return nil
	}

	size := binary.BigEndian.Uint32(w.buffer[1:5])

	if l < int(size+6) { // Payload not complete
		return nil
	}

	if w.buffer[0] != HEADER || w.buffer[size+5] != FOOTER {
		panic("Invalid container boundaries")
	}

	payload := w.buffer[5 : 5+size]
	w.buffer = w.buffer[:6+size]

	return payload
}

func Pack(payload []byte) []byte {
	size := 6 + len(payload)
	p := make([]byte, size)
	p[0] = HEADER
	p[size-1] = FOOTER
	writeUint32(p, uint32(len(payload)), 1)
	// TODO : This looks slow
	for i, b := range payload {
		p[5+i] = b
	}
	return p
}

func NewUnpacker() func([]byte) []byte {
	w := wire{make([]byte, 0)}
	return w.Unpack
}

func writeUint32(dst []byte, value uint32, offset uint32) {
	_ = dst[offset+3]
	dst[offset] = byte(value >> 24)
	dst[offset+1] = byte(value >> 16)
	dst[offset+2] = byte(value >> 8)
	dst[offset+3] = byte(value)
}
