package gate

import (
	"github.com/tieing/lemon/network/inet"
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
	"google.golang.org/protobuf/proto"
)

func (g *Gate) Push(_ rpc.Server, msg *rpcmsg.RPCMessage) {
	if conn := g.ws.GetConnMgr().GetConn(msg.CID); conn != nil {
		conn.Push(msg.BUF)
	}
}

func (g *Gate) PushAll(_ rpc.Server, msg *rpcmsg.RPCMessage) {
	g.ws.GetConnMgr().Each(func(conn inet.Conn) bool {
		conn.Push(msg.BUF)
		return true
	})
}

func (g *Gate) PushGroup(_ rpc.Server, msg *rpcmsg.RPCMessage) {
	data := &rpcmsg.PushGroupMsg{}
	err := proto.Unmarshal(msg.BUF, data)
	if err != nil {
		println(err.Error())
		return
	}
	for _, s := range data.Member {
		if conn := g.ws.GetConnMgr().GetConn(s); conn != nil {
			conn.Push(data.Data)
		}
	}
}

func (g *Gate) Bind(s rpc.Server, msg *rpcmsg.RPCMessage) {
	if conn := g.ws.GetConnMgr().GetConn(msg.CID); conn != nil {
		conn.Bind(int(msg.UID))
		s.Response([]byte("1"))
		return
	}
	s.Response(nil)
}
