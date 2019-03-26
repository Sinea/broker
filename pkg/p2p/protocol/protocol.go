package protocol

import "errors"

const ProtocolMagicNumber = 0xCE

type State uint

const (
	Idle State = 1 << iota
)

type Protocol interface {
	Write(flags MessageFlag, body []byte)
	Read() (message *Message, err error)
}

type proto struct {
	buffer []byte
}

func (p *proto) Write(flags MessageFlag, body []byte) {

}

func (p *proto) Read() (message *Message, err error) {
	if p.buffer[0] != ProtocolMagicNumber {
		return nil, errors.New(`invalid message header`)
	}

	return &Message{}, nil
}

func New(data []byte) Protocol {
	return &proto{
		buffer: data,
	}
}
