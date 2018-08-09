package io

import (
	"net"
	"time"
	"sync"
)

type Writer interface {
	Write(buffer []byte)
}

type socketWriter struct {
	connection net.Conn
	queue      [][]byte
	isWriting  bool
	lock       sync.Mutex
}

func (s *socketWriter) Write(buffer []byte) {
	s.queue = append(s.queue, buffer)
	s.lock.Lock()
	if !s.isWriting {
		s.isWriting = true
		go write(s, s.connection)
	}
	s.lock.Unlock()
}

func write(writer *socketWriter, conn net.Conn) {
	for {
		if len(writer.queue) == 0 {
			break
		}

		if err := conn.SetWriteDeadline(time.Now().Add(time.Millisecond)); err != nil {
			panic(err)
		}

		buffer := writer.queue[0]
		writer.queue = writer.queue[1:]
		n, err := conn.Write(buffer)

		if err != nil {
			if e, ok := err.(net.Error); ok {
				if !e.Temporary() && !e.Timeout() {
					panic(e)
				}
			} else {
				panic(err)
			}
		}

		if n < len(buffer) { // Wrote less bytes than there are in the current buffer
			writer.queue = append([][]byte{buffer[n:]}, writer.queue...)
		}
		time.Sleep(time.Millisecond)
	}

	writer.isWriting = false
}

func NewWriter(conn net.Conn) Writer {
	return &socketWriter{
		connection: conn,
		queue:      make([][]byte, 0),
		isWriting:  false,
		lock:       sync.Mutex{},
	}
}
