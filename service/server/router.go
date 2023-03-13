package server

import (
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
)

func (g *Server) RegisterServerMsgHandler(msgID string, handle func(rpc.Server, *rpcmsg.RPCMessage)) {
	g.processor.RegisterServerMsgHandler(msgID, handle)
}

func (g *Server) registerDefaultHandler() {
	g.RegisterServerMsgHandler("__OnSessionDisconnect", g.sessionDisconnect)

}
