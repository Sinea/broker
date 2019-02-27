package main

import (
	"github.com/Sinea/broker/pkg/broker"
)

func main() {
	b := broker.New(3)
	for _, t := range string("abcdefghijklmno1234567890!@#$%^&*()?><:{") {
		b.Write(string(t), []byte(string(t)))

	}
}
