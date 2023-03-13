package processor_ipml

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tieing/lemon/network/inet"
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
	"github.com/tieing/lemon/service/gate/gatemsg"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	"reflect"
)

type handleType func(rpc.Server, *rpcmsg.RPCMessage)

type RouterInfo struct {
	MsgId  int32
	Method reflect.Value
	Param  reflect.Type
}

type ProcessorImpl struct {
	serviceName     string
	codec           inet.Codec
	serverHandlers  map[string]handleType
	sessionHandlers map[int32]*RouterInfo
	sessMsg2ID      map[reflect.Type]int32
}

func NewProcessor(serviceName string, codec inet.Codec) *ProcessorImpl {
	return &ProcessorImpl{
		serviceName:     serviceName,
		codec:           codec,
		serverHandlers:  map[string]handleType{},
		sessionHandlers: map[int32]*RouterInfo{},
		sessMsg2ID:      map[reflect.Type]int32{},
	}
}

func (p *ProcessorImpl) OnServerMessage(svr rpc.Server, msg *rpcmsg.RPCMessage) {
	if fn, ok := p.serverHandlers[msg.HID]; ok {
		fn(svr, msg)
	} else {
		fmt.Printf("调用未注册的Server函数: %s", msg.HID)
	}
}

// 网关转发的客户端消息
// 在协程中执行,若要使消息有序列,需要重新实现这个方法,将消息写入队列等方式
func (p *ProcessorImpl) OnSessionMessage(sess rpc.Session, msg *rpcmsg.RPCMessage) {
	// session消息
	if info := p.sessionHandlers[msg.MID]; info != nil {
		// 解析
		var in = []reflect.Value{reflect.ValueOf(sess)}

		if info.Param != nil {
			p2 := reflect.New(info.Param.Elem())
			err := p.codec.Unmarshal(msg.BUF, p2.Interface())
			if err != nil {
				log.Err(err).Msgf("数据解析错误")
				return
			}
			in = append(in, p2)
			// -----log----
			log.Debug().Msgf("收到消息:MsgID:%d, CID:%d, UID:%d,  data:%s", msg.MID, sess.ID(), sess.UID(), p2.Interface().(protoiface.MessageV1).String())
			// -----log----
		}

		info.Method.Call(in)

	} else {
		log.Warn().Msgf("msg NodeId is not register! id = %d", msg.MID)
	}
}

func (p *ProcessorImpl) RegisterServerMsgHandler(msgID string, f func(rpc.Server, *rpcmsg.RPCMessage)) {
	p.serverHandlers[msgID] = f
}

func (p *ProcessorImpl) RegisterSessionMsgHandler(msgID int32, param any, f any) {
	info := RouterInfo{}
	info.MsgId = msgID
	info.Method = reflect.ValueOf(f)
	if param != nil {
		info.Param = reflect.TypeOf(param)
	}
	p.sessionHandlers[msgID] = &info
}

func (p *ProcessorImpl) RegisterResponseMsg(msgID int32, response any) {
	if _, ok := p.sessMsg2ID[reflect.TypeOf(response)]; ok {
		panic("response message registered!! " + reflect.TypeOf(response).Name())
	}
	p.sessMsg2ID[reflect.TypeOf(response)] = msgID
}

func (p *ProcessorImpl) MakeResponse(response any, code ...int32) []byte {
	var codeValue int32
	if len(code) > 0 {
		codeValue = code[0]
	}

	msgID, ok := p.sessMsg2ID[reflect.TypeOf(response)]
	if !ok {
		println("请注册响应消息!!! " + reflect.TypeOf(response).Name())
		return nil
	}
	var respByte []byte
	var err error
	if response != nil {
		respByte, err = p.codec.Marshal(response)
		if err != nil {
			println(err.Error())
			return nil
		}
		// -----log----
		log.Debug().Msgf("发送消息:CID:%d,data:%s", msgID, response.(protoiface.MessageV1).String())
		// -----log----
	}

	respByte, _ = proto.Marshal(&gatemsg.NetworkCMDMessage{
		Service: p.serviceName,
		Mid:     msgID,
		Buf:     respByte,
		Code:    codeValue,
		Cps:     false,
		Ect:     false,
	})
	return respByte
}
