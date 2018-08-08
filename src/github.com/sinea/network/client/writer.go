package client

import "github.com/sinea/network/io"

type MessageWriter interface {
	Write(serializable Message)
}

type messageWriter struct {
	next io.Writer
}

func (m *messageWriter) Write(message Message) {
	buffer := []byte{message.Kind(), message.Flags()}
	buffer = append(buffer, message.Body()...)

	m.next.Write(buffer)
}

func NewMessageWriter(next io.Writer) MessageWriter {
	return &messageWriter{next: next}
}
