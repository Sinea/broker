package main

import "broker/pkg/broker"

func main() {
	b := broker.New()
	b.Write("messages", []byte("john"))
	b.Write("messagez", []byte("john"))
	b.Write("messag", []byte("john"))
}
