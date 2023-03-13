package rpc_impl

import (
	"fmt"
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
	"google.golang.org/protobuf/proto"
)

type server struct {
	rpcClient      rpc.RPC
	remoteServerID string
	seqID          int32
}

func NewServer(client rpc.RPC, serverID string, seqID int32) rpc.Server {
	return &server{
		rpcClient:      client,
		remoteServerID: serverID,
		seqID:          seqID,
	}
}

func (s *server) Call(handlerID string, msg []byte) ([]byte, error) {
	r, err := s.rpcClient.Call(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_Request,
		RID:  s.remoteServerID,
		HID:  handlerID,
		BUF:  msg,
	})
	if err != nil {
		return nil, err
	}
	return r.BUF, nil
}

func (s *server) Notify(handlerID string, msg []byte) {
	s.rpcClient.Push(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_Notify,
		RID:  s.remoteServerID,
		HID:  handlerID,
		BUF:  msg,
	})
}
func (s *server) Response(msg []byte) {
	s.rpcClient.Push(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_Response,
		RID:  fmt.Sprintf(responseTopic, s.remoteServerID),
		SEQ:  s.seqID,
		BUF:  msg,
	})
}

func (s *server) ID() string {
	return s.remoteServerID
}

// ============ session ==========================================

type SessionImpl struct {
	uid       int
	cid       string // 连接ID
	rpcClient *RPCImpl
	gateID    string // 用户所在网关ID
}

func NewSession(uid int, cid, gid string, rpc *RPCImpl) *SessionImpl {
	return &SessionImpl{
		uid:       uid,
		cid:       cid,
		rpcClient: rpc,
		gateID:    gid,
	}
}

// ID 获取连接ID
func (sess *SessionImpl) ID() string {
	return sess.cid
}

// UID 获取用户ID
func (sess *SessionImpl) UID() int {
	return sess.uid
}

// UID 获取网关
func (sess *SessionImpl) GateID() string {
	return sess.gateID
}

// Push 发送消息（异步）
func (sess *SessionImpl) Push(data []byte) {
	if data == nil {
		return
	}
	sess.rpcClient.Push(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_S2C,
		RID:  sess.gateID,
		HID:  "__Push",
		CID:  sess.cid,
		BUF:  data,
	})
}

// PushAll 全网关推送
func (sess *SessionImpl) PushAll(data []byte) {
	if data == nil {
		return
	}
	sess.rpcClient.Push(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_S2C,
		RID:  sess.gateID,
		HID:  "__PushAll",
		CID:  "",
		BUF:  data,
	})
}

// PushAll 分组推送
func (sess *SessionImpl) PushGroup(data []byte, members []string) {
	if len(data) == 0 {
		return
	}
	data, _ = proto.Marshal(&rpcmsg.PushGroupMsg{
		Data:   data,
		Member: members,
	})
	sess.rpcClient.Push(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_S2C,
		RID:  sess.gateID,
		HID:  "__PushGroup",
		CID:  "",
		BUF:  data,
	})
}

// PushAll 分组推送
func (sess *SessionImpl) Bind(uid int) (bool, error) {
	sess.uid = uid
	r, err := sess.rpcClient.Call(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_Request,
		RID:  sess.gateID,
		HID:  "__Bind",
		CID:  sess.cid,
		UID:  int32(uid),
		BUF:  nil,
	})
	if err != nil {
		return false, err
	}
	return string(r.BUF) == "1", nil

}
