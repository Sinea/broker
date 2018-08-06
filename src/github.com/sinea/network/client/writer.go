package client

import "github.com/sinea/network/io"

type Message interface {
	Kind() uint8
	Flags() uint8
	Body() []byte
}

type message struct {
	kind  uint8
	flags uint8
	body  []byte
}

func (m *message) Kind() uint8 {
	return m.kind
}

func (m *message) Flags() uint8 {
	return m.flags
}

func (m *message) Body() []byte {
	return m.body
}

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

func NewMessage(kind uint8, body []byte) Message {
	return &message{kind, 0, body}
}
