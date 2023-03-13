package router

import (
	"github.com/rs/zerolog/log"
	"github.com/tieing/lemon/network/inet"

	"reflect"
	"sync"
)

type Router struct {
	routers sync.Map // map[int32]*RouterInfo
	respMap sync.Map // map[reflect.Type]int32
}

func NewRouter() *Router {
	return &Router{
		routers: sync.Map{},
		respMap: sync.Map{},
	}
}

func (r *Router) RegisterHandle(msgId int32, param interface{}, fn interface{}) {
	if _, ok := r.routers.Load(msgId); ok {
		log.Error().Msgf("msg NodeId is exist! id = %d", msgId)
		return
	}

	info := RouterInfo{}
	info.MsgId = msgId
	info.Method = reflect.ValueOf(fn)
	if param != nil {
		info.Param = reflect.TypeOf(param)
	}

	r.routers.Store(msgId, &info)

}

func (r *Router) RegisterMsg(msgId int32, msg interface{}) {
	if _, ok := r.respMap.Load(reflect.TypeOf(msg)); ok {
		log.Error().Msgf("msg NodeId is exist! id = %d", msgId)
		return
	}
	r.respMap.Store(reflect.TypeOf(msg), msgId)
}

func (r *Router) LoadRouterInfo(msgId int32) inet.IRouterInfo {
	if v, ok := r.routers.Load(msgId); ok {
		return v.(*RouterInfo)
	}
	return nil
}

func (r *Router) LoadMsgId(responseMsg interface{}) int32 {
	if v, ok := r.respMap.Load(reflect.TypeOf(responseMsg)); ok {
		return v.(int32)
	}
	return -1
}
