package router

import "reflect"

type RouterInfo struct {
	MsgId  int32
	Method reflect.Value
	Param  reflect.Type
}

func (info *RouterInfo) GetMsgID() int32                { return info.MsgId }
func (info *RouterInfo) GetMethod() reflect.Value       { return info.Method }
func (info *RouterInfo) GetParam() reflect.Type         { return info.Param }
func (info *RouterInfo) SetMsgID(msgID int32)           { info.MsgId = msgID }
func (info *RouterInfo) SetMethod(method reflect.Value) { info.Method = method }
func (info *RouterInfo) SetParam(param reflect.Type)    { info.Param = param }
