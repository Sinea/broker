package protocol

import "testing"

func TestProto_Write2(t *testing.T) {
	p := New([]byte{ProtocolMagicNumber + 1})
	m, err := p.Read()

	if m != nil {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestProto_Read_Success(t *testing.T) {
	p := New([]byte{ProtocolMagicNumber})
	m, err := p.Read()

	if m == nil {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestProto_Write(t *testing.T) {
	p := New([]byte{ProtocolMagicNumber})
	p.Write(System, []byte{})
}
