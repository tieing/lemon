package packet

type IMessage interface {
	GetMsgId() int32
	GetBody() []byte
	GetBodyLen() int
}

type Message struct {
	Route  int32  // 路由ID
	Buffer []byte // 消息内容
}

func (mp *Message) GetMsgId() int32 {
	return mp.Route
}

func (mp *Message) GetBody() []byte {
	return mp.Buffer
}

func (mp *Message) GetBodyLen() int {
	return len(mp.Buffer)
}
