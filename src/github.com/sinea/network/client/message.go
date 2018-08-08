package client

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

func NewMessage(kind uint8, body []byte) Message {
	return &message{kind, 0, body}
}
