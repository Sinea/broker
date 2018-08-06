package wire

import (
	"github.com/sinea/network/io"
)

type wireWriter struct {
	next io.Writer
}

func (w *wireWriter) Write(buffer []byte) {
	b := []byte{0xB0, uint8(len(buffer))}
	b = append(b, buffer...)
	b = append(b, 0x0B)

	w.next.Write(b)
}

func NewWriter(next io.Writer) io.Writer {
	return &wireWriter{next: next}
}
