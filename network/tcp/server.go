package tcp

import (
	"github.com/tieing/lemon/network/inet"
	"net"
)

type server struct {
	serverID string
	opts     *serverOptions // 配置
	listener net.Listener   // 监听器
	connMgr  *connMgr       // 连接管理器

	startHandler      inet.StartHandler      // 服务器启动hook函数
	stopHandler       inet.CloseHandler      // 服务器关闭hook函数
	connectHandler    inet.ConnectHandler    // 连接打开hook函数
	disconnectHandler inet.DisconnectHandler // 连接关闭hook函数
	receiveHandler    inet.ReceiveHandler    // 接收消息hook函数
	mateStorage       inet.MetaStorage       // 连接session数据缓存处理函数

}

var _ inet.Server = &server{}

func NewServer(opts ...ServerOption) inet.SocketServer {
	o := defaultServerOptions()
	for _, opt := range opts {
		opt(o)
	}

	s := &server{}
	s.opts = o
	s.connMgr = newConnMgr(s)

	return s
}

// Addr 监听地址
func (s *server) Addr() string {
	return s.opts.addr
}

// Start 启动服务器
func (s *server) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	if s.startHandler != nil {
		s.startHandler()
	}

	return s.serve()
}

// Stop 关闭服务器
func (s *server) Stop() error {
	if err := s.listener.Close(); err != nil {
		return err
	}

	s.connMgr.close()

	return nil
}

// Protocol 协议
func (s *server) Protocol() string {
	return "tcp"
}

// OnStart 监听服务器启动
func (s *server) OnStart(handler inet.StartHandler) {
	s.startHandler = handler
}

// OnStop 监听服务器关闭
func (s *server) OnStop(handler inet.CloseHandler) {
	s.stopHandler = handler
}

// OnConnect 监听连接打开
func (s *server) OnConnect(handler inet.ConnectHandler) {
	s.connectHandler = handler
}

// OnDisconnect 监听连接关闭
func (s *server) OnDisconnect(handler inet.DisconnectHandler) {
	s.disconnectHandler = handler
}

// OnReceive 监听接收到消息
func (s *server) OnReceive(handler inet.ReceiveHandler) {
	s.receiveHandler = handler
}

func (s *server) GetConnMgr() inet.ConnMgr {
	return s.connMgr
}

func (s *server) SetCacheStorage(storage inet.MetaStorage) {
	s.mateStorage = storage
}

// 初始化TCP服务器
func (s *server) init() error {
	addr, err := net.ResolveTCPAddr("tcp", s.opts.addr)
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return err
	}

	s.listener = ln

	return nil
}

// 等待连接
func (s *server) serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				continue
			}
		}
		s.connMgr.allocate(s, conn)
	}
}
