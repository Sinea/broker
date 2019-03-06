package main

import (
	"github.com/Sinea/broker/pkg/p2p"
	"time"
)

func main() {
	mesh := p2p.New()
	go mesh.Listen("0.0.0.0:1111")

	time.Sleep(time.Minute)
}
