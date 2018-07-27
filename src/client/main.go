package main

import (
	"net"
	"log"
	"network"
	"time"
)

func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", "0.0.0.0:3333")
	if err != nil {
		panic("Shit")
	}

	return conn
}

func HealthCheck(client network.Client, interval time.Duration, timeout time.Duration) {
	for {
		log.Println("Ping")
		select {
		case reply := <-client.Request(network.Message{Kind: network.PING}):
			if reply.Kind != network.PONG {
				panic("Wring ping reply")
			}
			log.Println("Pong")
			break
		case <-time.After(time.Second):
			panic("Ping timed out")
			break
		}

		time.Sleep(interval)
	}
}

func main() {
	conn := connectToServer()
	defer conn.Close()
	client := network.NewClient(conn)
	HealthCheck(client, time.Second, 500*time.Millisecond)
}
