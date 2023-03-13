package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/tieing/lemon/network/inet"
	"net"
	"net/http"
)

type server struct {
	opts     *serverOptions // 配置
	listener net.Listener   // 监听器
	connMgr  *connMgr       // 连接管理器
	serverID string

	startHandler      inet.StartHandler        // 服务器启动hook函数
	stopHandler       inet.CloseHandler        // 服务器关闭hook函数
	connectHandler    inet.ConnectHandler      // 连接打开hook函数
	disconnectHandler inet.DisconnectHandler   // 连接关闭hook函数
	receiveHandler    inet.ReceiveHandler      // 接收消息hook函数
	mateStorage       inet.MetaStorage         // 连接session数据缓存处理函数
	beforeUpgrade     inet.BeforeUpgradeHandle // http 升级ws之前调用
}

var _ inet.Server = &server{}

func NewServer(serverID string, opts ...ServerOption) inet.SocketServer {
	o := defaultServerOptions()
	for _, opt := range opts {
		opt(o)
	}

	s := &server{}
	s.opts = o
	s.serverID = serverID
	s.connMgr = newConnMgr(s)

	return s
}

// Addr 监听地址
func (s *server) Addr() string {
	return s.opts.addr
}

// Protocol 协议
func (s *server) Protocol() string {
	return "websocket"
}

// Start 启动服务器
func (s *server) Start() error {
	if err := s.init(); err != nil {
		return err
	}
	if s.receiveHandler == nil {
		fmt.Println("not set receiveHandler,now use default handle, pleas use server.OnReceive(func) setting it!")
		s.receiveHandler = func(conn inet.Conn, msg []byte) {
			fmt.Printf("connID:%d,BindID:%d,msg:%s", conn.ID(), conn.UID(), msg)
		}
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

// 初始化服务器
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

// 启动服务器
func (s *server) serve() error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin:       s.opts.checkOrigin,
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(s.opts.path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", 405)
			return
		}

		// 升级到ws之前调用,若返回false, 则这个http协议不能升级为web socket
		if s.beforeUpgrade != nil && !s.beforeUpgrade(w, r) {
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Error().Msgf("websocket upgrade error: %v", err)
			return
		}

		if int(s.connMgr.count) > s.opts.maxConnNum {
			//return inet.ErrTooManyConnection
			_ = conn.Close()
		}
		s.connMgr.allocate(s, conn)
	})

	if s.opts.certFile != "" && s.opts.keyFile != "" {
		return http.ServeTLS(s.listener, serveMux, s.opts.certFile, s.opts.keyFile)
	}

	return http.Serve(s.listener, serveMux)
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

// 获取codec
//func (s *server) GetCodec() inet.Codec {
//	return s.opts.codec
//}
//
//// 获取路由
//func (s *server) GetRouter() inet.IRouter {
//	return s.opts.router
//}
//
//func (s *server) GetPacker() inet.Packer {
//	return s.opts.packer
//}
