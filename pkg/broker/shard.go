package broker

type shard struct {
	peers []peer
}

func (s *shard) Write(topic string, data []byte) {
	// Write local
}
