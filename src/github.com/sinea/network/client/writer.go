package client

import "github.com/sinea/network/io"

type MessageWriter interface {
	Write(serializable Message)
}

type messageWriter struct {
	next io.Writer
}

func (m *messageWriter) Write(serializable Message) {
	buffer := []byte{serializable.Kind(), serializable.Flags()}
	buffer = append(buffer, serializable.Body()...)

	m.next.Write(buffer)
}

func NewMessageWriter(next io.Writer) MessageWriter {
	return &messageWriter{next: next}
}
