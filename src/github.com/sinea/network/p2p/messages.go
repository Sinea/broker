package p2p

type tag struct {
	tag uint16
}

// Message sent by the peer that connects
type PeerAuth struct {
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
}

type AskRequest struct {
	*tag
	c Client
	Data string
}

func (r *AskRequest) Reply(reply interface{}) {
	reply.(*tag).tag = r.tag.tag
	r.c.Message(reply)
}

type AskReply struct {
	*tag
	Result string `json:"result"`
}

func NewAskReply(result string) AskReply {
	return AskReply{Result: result}
}

func NewAskRequest(what string) AskRequest {
	return AskRequest{&tag{0},nil, what}
}