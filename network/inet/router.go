package inet

import "reflect"

type IRouter interface {
	RegisterHandle(msgId int32, param interface{}, fn interface{})
	RegisterMsg(msgId int32, msg interface{})
	LoadRouterInfo(msgId int32) IRouterInfo
	LoadMsgId(responseMsg interface{}) int32
}

type IRouterInfo interface {
	GetMsgID() int32
	GetMethod() reflect.Value
	GetParam() reflect.Type
	SetMsgID(msgID int32)
	SetMethod(method reflect.Value)
	SetParam(param reflect.Type)
}
