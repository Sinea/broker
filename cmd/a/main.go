package main

import (
	"fmt"
	"github.com/Sinea/broker/pkg/broker"
	"log"
	"time"
)

var counter = 0

func listen(b broker.Broker, prefix, topic string) {
	for range b.Read(topic) {
		//fmt.Printf("%s: %s\n", prefix, string(t))
		counter++
	}
}

func start(b broker.Broker) {
	if err := b.Start("0.0.0.0:4444"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	b := broker.New(0)
	go listen(b, "on x", "x")
	go start(b)

	//time.Sleep(time.Second*5)

	//if err := b.Join("0.0.0.0:3333"); err != nil {
	//	log.Fatal(err)
	//}

	t := time.Now()
	for {
		b.Write("x", []byte("msg1"))
		//b.Write("z", []byte("msg2"))
		time.Sleep(time.Nanosecond)
		if time.Since(t) >= time.Second {
			fmt.Printf("Got %d messages in %s\n", counter, time.Since(t))
			t = time.Now()
			counter = 0
		}
	}

}
