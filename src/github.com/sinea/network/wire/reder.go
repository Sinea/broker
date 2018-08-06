package wire

import "github.com/sinea/network/io"

type wireHandler struct {
	buffer []byte
	next   io.Writer
}

func (w *wireHandler) Write(buffer []byte) {
	w.buffer = append(w.buffer, buffer...)

	if w.buffer[0] != 0xB0 {
		panic("Invalid header")
	}

	size := w.buffer[1]
	if uint8(len(w.buffer)) < size+2 {
		return
	}

	w.next.Write(w.buffer[2:size+2])
	w.buffer = w.buffer[:size+2]
}

func NewReader(writer io.Writer) io.Writer {
	return &wireHandler{make([]byte, 0), writer}
}
