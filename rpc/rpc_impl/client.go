package rpc_impl

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
	"sync/atomic"
	"time"
)

var ErrTimeOut = errors.New("rpc: timeout")
var ErrNoKnow = errors.New("rpc: unknow")

const CALL_TIMEOUT = 10 * time.Second

// 协程写入
func (r *RPCImpl) writeLoop() {
	var item any
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			if item = r.writeQue.Remove(); item != nil {
				r.send(item.(*rpcmsg.RPCMessage))
			} else {
				time.Sleep(time.Nanosecond)
			}
		}
	}
}

// 阻塞式
func (p *RPCImpl) Call(req *rpcmsg.RPCMessage) (resp *rpcmsg.RPCMessage, err error) {
	call := p.newCall()

	req.SEQ = call.seqID
	req.SID = p.serverID

	buf, _ := proto.Marshal(req)
	err = p.conn.Publish(req.RID, buf)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(p.ctx, CALL_TIMEOUT)
	select {
	case <-ctx.Done():
		cancel()
		p.getCallWithDel(call.seqID)
		return nil, ctx.Err()
	case data := <-call.Ch:
		cancel()
		p.getCallWithDel(call.seqID)
		return data, call.Err
	}
}

func (p *RPCImpl) send(msg *rpcmsg.RPCMessage) {
	buf, _ := proto.Marshal(msg)
	err := p.conn.Publish(msg.RID, buf)
	if err != nil {
		panic(err)
	}
}

// 异步发送
func (p *RPCImpl) Push(msg *rpcmsg.RPCMessage) {
	msg.SID = p.serverID
	p.writeQue.Add(msg)
}

// ------------------------------------------

func (p *RPCImpl) newSeqID() int32 {
	return atomic.AddInt32(&p.seqID, 1)
}

type Call struct {
	seqID int32
	Ch    chan *rpcmsg.RPCMessage
	Err   error
}

func (p *RPCImpl) newCall() *Call {
	seqID := p.newSeqID()
	call := &Call{
		seqID: seqID,
		Ch:    make(chan *rpcmsg.RPCMessage),
	}
	p.callSet.Store(seqID, call)
	return call
}

func (p *RPCImpl) getCallWithDel(seqID int32) (*Call, bool) {
	if v, ok := p.callSet.LoadAndDelete(seqID); ok {
		return v.(*Call), true
	}
	return nil, false
}

func (p *RPCImpl) handleResponse(msg *rpcmsg.RPCMessage) error {
	call, ok := p.getCallWithDel(msg.SEQ)
	if !ok {
		return errors.New("seqID not existed")
	}
	call.Err = nil
	call.Ch <- msg

	return call.Err
}
