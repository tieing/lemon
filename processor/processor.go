package processor

import (
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
)

type Processor interface {
	OnServerMessage(svr rpc.Server, msg *rpcmsg.RPCMessage)
	OnSessionMessage(sess rpc.Session, msg *rpcmsg.RPCMessage)
	RegisterServerMsgHandler(msgID string, f func(rpc.Server, *rpcmsg.RPCMessage))
}

type DefaultProcessor interface {
	Processor
	RegisterSessionMsgHandler(msgID int32, param any, f any)
	MakeResponse(response any, code ...int32) []byte
	RegisterResponseMsg(msgID int32, response any)
}
