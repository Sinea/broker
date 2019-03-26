package protocol

type MessageFlag uint8

func (f MessageFlag) Has(flag MessageFlag) bool {
	return (f & flag) != 0
}

const System MessageFlag = 1 << iota
