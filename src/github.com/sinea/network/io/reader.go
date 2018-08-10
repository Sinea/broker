package io

import (
	"net"
	"time"
)

func NewReader(conn net.Conn, writer Writer) {
	go read(conn, writer)
}

func read(conn net.Conn, writer Writer) {
	b := make([]byte, 1024)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			panic(err)
		}

		n, err := conn.Read(b)

		if err != nil {
			if e, ok := err.(net.Error); ok {
				if !e.Timeout() {
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
