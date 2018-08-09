package io

import (
	"net"
	"time"
)

func NewReader(conn net.Conn, writer Writer) {
	go read(conn, writer)
}

func read(conn net.Conn, writer Writer) {
	for {
		if err := conn.SetReadDeadline(time.Now().Add(33 * time.Millisecond)); err != nil {
			panic(err)
		}

		b := make([]byte, 1024)
		n, err := conn.Read(b)

		if err != nil {
			if e, ok := err.(net.Error); ok {
				if !e.Temporary() && !e.Timeout() {
					panic(e)
				}
			} else {
				panic(err)
			}
		}

		if n != 0 {
			writer.Write(b[:n])
		}
	}
}
