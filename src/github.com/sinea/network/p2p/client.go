package p2p

type MessageSender interface {
	Message(message interface{})
}

type MessageReceiver interface {
	Messages() <-chan interface{}
}

type RequestSender interface {
	Request(request interface{}, options RequestOptions) <-chan interface{}
}

type RequestOptions struct {
	Timeout uint32
}

type Request interface {
	Reply(reply interface{})
}

type Client interface {
	MessageSender
	MessageReceiver
	RequestSender

	Connect(address string) error
}

func NewClient() Client {
	return nil
}
