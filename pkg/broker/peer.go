package broker

type peer struct {

}

func (p *peer) Write(topic string, data []byte) {
	// Send by socket
}

