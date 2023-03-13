package rpc_impl

import (
	"context"
	"fmt"
	"github.com/tieing/lemon/processor"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
	"github.com/tieing/lemon/tools/queue"
	"github.com/gammazero/workerpool"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"sync"
)

const (
	responseTopic = "%s:response"
)

type RPCImpl struct {
	ctx        context.Context
	cancel     context.CancelFunc
	serverID   string
	conn       *nats.Conn
	sub        []*nats.Subscription
	writeQue   *queue.Queue
	workerPool *workerpool.WorkerPool
	processor  processor.Processor

	// ------
	seqID   int32
	callSet sync.Map
}

func NewRPC(ctx context.Context, serverID string, processor processor.Processor, conn *nats.Conn) *RPCImpl {
	ctx, cancel := context.WithCancel(ctx)
	rpc := &RPCImpl{
		ctx:        ctx,
		cancel:     cancel,
		serverID:   serverID,
		conn:       conn,
		writeQue:   queue.New(),
		workerPool: workerpool.New(4096),
		processor:  processor,
	}
	return rpc
}

// 订阅
func (r *RPCImpl) subscribe() (err error) {
	// 订阅请求消息
	sub, err := r.conn.Subscribe(r.serverID, func(msg *nats.Msg) {
		r.workerPool.Submit(func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("%v", e)
				}
			}()
			r.msgHandler(msg.Data)
		})
	})
	if err != nil {
		return err
	}
	r.sub = append(r.sub, sub)

	// 订阅响应消息
	// 解决 大量请求后响应被阻塞,但是前面的请求却一直等待响应结果,最终响应超时的问题
	sub, err = r.conn.Subscribe(fmt.Sprintf(responseTopic, r.serverID), func(msg *nats.Msg) {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("%v", e)
			}
		}()
		r.msgHandler(msg.Data)
	})
	if err != nil {
		return err
	}
	r.sub = append(r.sub, sub)
	return
}

func (r *RPCImpl) Run() (err error) {
	go r.writeLoop()
	return r.subscribe()
}

func (r *RPCImpl) Close() {
	for _, subscription := range r.sub {
		subscription.Unsubscribe()
	}

	r.cancel()
}

func (r *RPCImpl) msgHandler(buf []byte) {
	msg := new(rpcmsg.RPCMessage)
	err := proto.Unmarshal(buf, msg)
	if err != nil {
		return
	}

	switch msg.Type {
	case rpcmsg.Type_S2C, rpcmsg.Type_Request, rpcmsg.Type_Notify:
		r.processor.OnServerMessage(NewServer(r, msg.SID, msg.SEQ), msg)

	case rpcmsg.Type_Response:
		err = r.handleResponse(msg)

	case rpcmsg.Type_C2S:
		r.processor.OnSessionMessage(NewSession(int(msg.UID), msg.CID, msg.SID, r), msg)
	default:
		return
	}

	if err != nil {
		println(err.Error())
	}
}
