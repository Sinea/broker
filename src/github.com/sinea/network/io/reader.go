package io

import (
	"net"
	"time"
)

func NewReader(conn net.Conn, writer Writer) {
	go func(c net.Conn) {
		for {
			c.SetReadDeadline(time.Now().Add(time.Millisecond))
			b := make([]byte, 1024)
			n, err := c.Read(b)
			if err != nil {
				if e, ok := err.(net.Error); ok && e.Timeout() {
					
				} else {
					panic(err)
				}
			} else {
				writer.Write(b[:n])
			}
			time.Sleep(time.Millisecond)
		}
	}(conn)
}
