package main

import (
	"bufio"
	"fmt"
	"github.com/Sinea/broker/pkg/p2p"
	"log"
	"os"
)

func dump(mesh p2p.Mesh) {
	for message := range mesh.Read() {
		fmt.Printf("Received '%s' from %d\n", string(message.Data), message.From)
	}
}

func main() {
	mesh := p2p.New(9)
	go mesh.Listen("0.0.0.0:3333")
	go dump(mesh)
	mesh.Join("0.0.0.0:2222")

	reader := bufio.NewReader(os.Stdin)

	for {
		bytes, _ := reader.ReadBytes('\n')
		log.Printf("Will send data %s", string(bytes))
		if peer, _ := mesh.Peer(2); peer != nil {
			peer.Send(bytes)
		} else {
			log.Println("Nil peer")
		}
	}
}
