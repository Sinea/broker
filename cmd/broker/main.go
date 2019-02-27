package main

import (
	"fmt"
	"github.com/Sinea/broker/pkg/broker"
	"time"
)

func listen(b broker.Broker, prefix, topic string) {
	for t := range b.Read(topic) {
		fmt.Println(prefix, t)
	}
}

func main() {
	b := broker.New(3)

	for i := 0; i < 3; i++ {
		go listen(b, fmt.Sprintf("%d.", i), "a")
	}

	for {
		b.Write("a", []byte(string("a")))
		time.Sleep(time.Second)
	}
}
