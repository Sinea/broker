package main

import (
	"github.com/Sinea/broker/pkg/p2p"
)

func main() {
	mesh := p2p.New()
	go mesh.Listen("0.0.0.0:2222")
	mesh.Join("0.0.0.0:1111")
}
