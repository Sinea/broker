package client

import "github.com/sinea/network/io"

type Reader interface {
	io.Writer
	Messages() <-chan Message
}

type reader struct {
	messages chan Message
}

func (r *reader) Write(buffer []byte) {
	r.messages <- NewMessage(buffer[0], buffer[2:])
}

func (r *reader) Messages() <-chan Message {
	return r.messages
}

func NewReader() Reader {
	return &reader{
		make(chan Message),
	}
}
