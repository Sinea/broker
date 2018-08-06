package io

import (
	"net"
	"log"
	"time"
)

type Writer interface {
	Write(buffer []byte)
}

type socketWriter struct {
	connection net.Conn
	queue      [][]byte
}

func (s *socketWriter) Write(buffer []byte) {
	log.Printf("Push : %X", buffer)
	s.queue = append(s.queue, buffer)
}

func NewWriter(conn net.Conn) Writer {
	s := socketWriter{connection: conn, queue: make([][]byte, 0)}

	go func(ss *socketWriter, cc net.Conn) {
		for {
			if len(ss.queue) != 0 {
				cc.SetWriteDeadline(time.Now().Add(time.Millisecond))
				b := ss.queue[0]
				ss.queue = ss.queue[1:]
				log.Printf("Write: %X", b)
				n, err := cc.Write(b)
				if err != nil {
					if err, ok := err.(net.Error); ok && err.Timeout() {
						log.Println("Just a timeout")
					} else {
						log.Fatalf("Strage error: %s", err)
					}
				}
				if n < len(b) {
					ss.queue = append([][]byte{b[n:]}, ss.queue...)
				}
			}
			time.Sleep(time.Millisecond)
		}
	}(&s, conn)

	return &s
}
