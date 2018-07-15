package main

import (
	"net"
	"log"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:3333")

	if err != nil {
		log.Println("Problem connecting")
	}

	for {
		conn.Write([]byte("Hello world"))
		time.Sleep(time.Second)
	}
}
