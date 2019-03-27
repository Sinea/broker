package p2p

import (
	"encoding/json"
	"log"
)

func bytes(message interface{}) []byte {
	b, err := json.Marshal(message)

	if err != nil {
		log.Fatal(err)
	}

	return b
}
