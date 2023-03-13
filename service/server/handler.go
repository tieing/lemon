package server

import (
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
)

func (g *Server) sessionDisconnect(_ rpc.Server, msg *rpcmsg.RPCMessage) {
	if g.onSessionDisconnect != nil {
		g.onSessionDisconnect(msg.SID, msg.CID, int(msg.UID))
	}
}
