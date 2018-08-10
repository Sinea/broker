package client

import (
	"testing"
	"log"
	"github.com/sinea/network/wire"
)

type mockWriter struct {
}

func (*mockWriter) Write(buffer []byte) {
	log.Printf("Packed message: %X", buffer)
}

func TestMessageWriter(t *testing.T) {
	mock := &mockWriter{}
	w := wire.NewWriter(mock)
	writer := NewMessageWriter(w)
	writer.Write(Hello{"xyz", 1})
}
