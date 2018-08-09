package wire

import (
	"github.com/sinea/network/io"
)

type wireHandler struct {
	buffer []byte
	next   io.Writer
}

func (w *wireHandler) Write(buffer []byte) {
	w.buffer = append(w.buffer, buffer...)

	// We need at least 2 bytes in the buffer (header + size)
	if len(w.buffer) < 2 {
		return
	}

	if w.buffer[0] != 0xB0 {
		panic("invalid container boundary")
	}

	size := w.buffer[1]
	if uint8(len(w.buffer)) < size+3 {
		return
	}

	if w.buffer[size+2] != 0x0B {
		panic("invalid container boundary")
	}

	payload := w.buffer[2 : size+2]
	w.buffer = w.buffer[size+3:]
	w.next.Write(payload)
}

func NewReader(writer io.Writer) io.Writer {
	return &wireHandler{make([]byte, 0), writer}
}
