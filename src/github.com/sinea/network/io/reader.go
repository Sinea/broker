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
				if err, ok := err.(net.Error); ok && err.Timeout() {
					//log.Println("Just a timeout error")
				}
			} else {
				writer.Write(b[:n])
			}
			time.Sleep(time.Millisecond)
		}
	}(conn)
}
