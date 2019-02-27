package main

import (
	"fmt"
	"github.com/Sinea/broker/pkg/broker"
	"time"
)

func listen(b broker.Broker, topic string) {
	for t := range b.Read(topic) {
		fmt.Println(t)
	}
}

func main() {
	b := broker.New(3)
	go listen(b, "a")
	go listen(b, "c")
	time.Sleep(time.Second)
	for _, t := range string("abcdefghijklmno1234567890!@#$%^&*()?><:{") {
		b.Write(string(t), []byte(string(t)))
	}

	time.Sleep(time.Second)
}
