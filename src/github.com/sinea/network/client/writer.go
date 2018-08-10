package client

import (
	"github.com/sinea/network/io"
	"encoding/json"
)

type MessageWriter interface {
	Write(m interface{})
}

type messageWriter struct {
	next io.Writer
}

func (m *messageWriter) Write(message interface{}) {
	b, err := json.Marshal(message)

	if err != nil {
		panic(err)
	}

	buffer := []byte{getID(message)}
	buffer = append(buffer, b...)

	m.next.Write(buffer)
}

func NewMessageWriter(next io.Writer) MessageWriter {
	return &messageWriter{next: next}
}
