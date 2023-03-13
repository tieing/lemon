package gate

import (
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
)

func (g *Gate) RegisterServerMsgHandler(msgID string, handle func(rpc.Server, *rpcmsg.RPCMessage)) {
	g.processor.RegisterServerMsgHandler(msgID, handle)
}

func (g *Gate) registerDefaultHandler() {
	g.RegisterServerMsgHandler("__Push", g.Push)
	g.RegisterServerMsgHandler("__PushGroup", g.PushGroup)
	g.RegisterServerMsgHandler("__PushAll", g.PushAll)
	g.RegisterServerMsgHandler("__Bind", g.Bind)

}
