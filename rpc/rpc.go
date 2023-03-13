package rpc

import "github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"

type RPC interface {
	Run() (err error)
	Close()
	Call(req *rpcmsg.RPCMessage) (resp *rpcmsg.RPCMessage, err error)
	Push(msg *rpcmsg.RPCMessage)
}

type Server interface {
	Call(string, []byte) ([]byte, error)
	Notify(method string, msg []byte)
	ID() string
	Response([]byte)
}

type Session interface {
	// ID 获取连接ID
	ID() string
	UID() int
	GateID() string
	// Push 发送消息
	Push(data []byte)
	// PushGroup 对指定的一个或者多个CID进行推送
	PushGroup(data []byte, group []string)
	// PushAll 全网关发送
	PushAll(data []byte)

	Bind(uid int) (bool, error)
}
