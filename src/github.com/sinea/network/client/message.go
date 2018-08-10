package client

type Message interface {
	Kind() uint8
	Flags() uint8
	Body() []byte
}

type message struct {
	kind  uint8
	flags uint8
	body  []byte
}

func (m *message) Kind() uint8 {
	return m.kind
}

func (m *message) Flags() uint8 {
	return m.flags
}

func (m *message) Body() []byte {
	return m.body
}

type Hello struct {
	Name string `json:"name"`
	Age  uint16 `json:"age"`
}

type Goodbye struct {
	From string
}

func getID(i interface{}) uint8 {
	switch i.(type) {
	case Hello:
		return 1
	case Goodbye:
		return 2
	default:
		panic("unknown type")
	}
}

func reverseID(i uint8) interface{} {
	switch i {
	case 1:
		return &Hello{}
	case 2:
		return &Goodbye{}
	default:
		panic("unknown id")
	}
}
