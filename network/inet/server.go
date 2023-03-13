package inet

import "net/http"

type (
	StartHandler        func()
	CloseHandler        func()
	ConnectHandler      func(conn Conn)
	DisconnectHandler   func(conn Conn)
	ReceiveHandler      func(conn Conn, msg []byte)
	BeforeUpgradeHandle func(w http.ResponseWriter, r *http.Request) bool
)

type Server interface {
	// Addr 监听地址
	Addr() string
	// Start 启动服务器
	Start() error
	// Stop 关闭服务器
	Stop() error
	// OnStart 监听服务器启动
	OnStart(handler StartHandler)
	// OnStop 监听服务器关闭
	OnStop(handler CloseHandler)
}

type SocketServer interface {
	// Addr 监听地址
	Addr() string
	// Start 启动服务器
	Start() error
	// Stop 关闭服务器
	Stop() error
	// OnStart 监听服务器启动
	OnStart(handler StartHandler)
	// OnStop 监听服务器关闭
	OnStop(handler CloseHandler)
	// OnConnect 监听连接打开
	OnConnect(handler ConnectHandler)
	// OnReceive 监听接收消息
	OnReceive(handler ReceiveHandler)
	// OnDisconnect 监听连接断开
	OnDisconnect(handler DisconnectHandler)
	//GetConnMgr 获取连接管理器
	GetConnMgr() ConnMgr
	// Protocol 协议
	Protocol() string
	// SetCacheStorage 设置缓存存储函数
	SetCacheStorage(storage MetaStorage)
}
