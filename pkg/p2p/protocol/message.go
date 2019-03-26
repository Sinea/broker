package protocol

type Message struct {
	From  uint8
	To    uint8
	Flags MessageFlag
	Data  []byte
}
