package wire

import (
	"testing"
	"github.com/sinea/network/io"
)

type mockWriter struct {

}

func (*mockWriter) Write(buffer []byte) {
	if buffer[0] != 0xB0 {
		panic("Bad header")
	}

	if buffer[len(buffer)-1] != 0x0B {
		panic("Bad footer")
	}
}

func TestWireWriter(t *testing.T) {
	var w io.Writer = &mockWriter{}
	w = NewWriter(w)
	w.Write([]byte("hello"))
}
