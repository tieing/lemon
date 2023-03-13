package inet

import (
	"errors"
	"net"
)

const (
	ConnOpened ConnState = iota + 1 // 连接打开
	ConnHanged                      // 连接挂起
	ConnClosed                      // 连接关闭
)

var (
	ErrConnectionHanged  = errors.New("connection is hanged")
	ErrConnectionClosed  = errors.New("connection is closed")
	ErrIllegalMsgType    = errors.New("illegal message type")
	ErrTooManyConnection = errors.New("too many connection")
)

type (
	ConnState int32

	Conn interface {
		// ID 获取连接ID
		ID() string
		// UID 获取用户ID
		UID() int
		// Bind 绑定用户ID
		Bind(uid int)
		// Unbind 解绑用户ID
		Unbind()
		// Send 发送消息（同步）
		Send(data []byte, msgType ...int) error
		// Push 发送消息（异步）
		Push(data []byte, msgType ...int) (err error)
		// PushAll 全网关发送(异步)
		PushAll(data []byte, smgType ...int) error
		// State 获取连接状态
		State() ConnState
		// Close 关闭连接
		Close(isForce ...bool) error
		// LocalIP 获取本地IP
		LocalIP() (string, error)
		// LocalAddr 获取本地地址
		LocalAddr() (net.Addr, error)
		// RemoteIP 获取远端IP
		RemoteIP() (string, error)
		// RemoteAddr 获取远端地址
		RemoteAddr() (net.Addr, error)

		MateCache
	}

	MateCache interface {
		GetMateData(string) any
		DelMateData(string)
		SetMetaData(k string, v any)
		GetMetaDataAll() map[string]any
	}

	// MetaStorage 缓存用户零时数据,用作断线重连等
	MetaStorage interface {
		Storage(uid any, m map[string]any) error
		Delete(uid any) error
		Select(uid any) (map[string]any, error)
	}
)
