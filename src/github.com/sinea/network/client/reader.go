package client

import (
	"github.com/sinea/network/io"
	"encoding/json"
)

type Reader interface {
	io.Writer
	Messages() <-chan interface{}
}

type reader struct {
	messages chan interface{}
}

func (r *reader) Write(buffer []byte) {
	var v = reverseID(buffer[0])
	json.Unmarshal(buffer[1:], &v)

	r.messages <- v
}

func (r *reader) Messages() <-chan interface{} {
	return r.messages
}

func NewReader() Reader {
	return &reader{
		make(chan interface{}),
	}
}
