package io

import (
	"net"
	"log"
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
	log.Printf("Push : %X", buffer)
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
		if len(writer.queue) != 0 {
			conn.SetWriteDeadline(time.Now().Add(time.Millisecond))
			b := writer.queue[0]
			writer.queue = writer.queue[1:]
			n, err := conn.Write(b)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					//log.Println("Just a timeout")
				} else {
					log.Fatalf("Strage error: %s", err)
					panic(err)
				}
			}
			if n < len(b) {
				writer.queue = append([][]byte{b[n:]}, writer.queue...)
			}
		}
		time.Sleep(time.Millisecond)
	}

	writer.isWriting = true
}

func NewWriter(conn net.Conn) Writer {
	return &socketWriter{
		connection: conn,
		queue:      make([][]byte, 0),
		isWriting:  false,
		lock:       sync.Mutex{},
	}
}
