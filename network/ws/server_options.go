package ws

import (
	"net/http"
	"time"
)

const (
	defaultServerAddr                   = ":3553"
	defaultServerPath                   = "/"
	defaultServerMaxMsgLen              = 1024
	defaultServerMaxConnNum             = 5000
	defaultServerCheckOrigin            = "*"
	defaultServerHeartbeatCheck         = false
	defaultServerHeartbeatCheckInterval = 1000
	defaultServerHandshakeTimeout       = 2000
)

type ServerOption func(o *serverOptions)

type CheckOriginFunc func(r *http.Request) bool

type serverOptions struct {
	addr                   string          // 监听地址
	maxMsgLen              int             // 最大消息长度（字节），默认1kb
	maxConnNum             int             // 最大连接数
	certFile               string          // 证书文件
	keyFile                string          // 秘钥文件
	path                   string          // 路径，默认为"/"
	checkOrigin            CheckOriginFunc // 跨域检测
	enableHeartbeatCheck   bool            // 是否启用心跳检测
	heartbeatCheckInterval time.Duration   // 心跳检测间隔时间，默认10s
	handshakeTimeout       time.Duration   // 握手超时时间，默认10s
}

func defaultServerOptions() *serverOptions {
	origins := []string{defaultServerCheckOrigin}
	checkOrigin := func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		for _, v := range origins {
			if v == defaultServerCheckOrigin || origin == v {
				return true
			}
		}
		return false
	}

	return &serverOptions{
		checkOrigin:            checkOrigin,
		addr:                   defaultServerAddr,
		maxMsgLen:              defaultServerMaxMsgLen,
		maxConnNum:             defaultServerMaxConnNum,
		path:                   defaultServerPath,
		keyFile:                "",
		certFile:               "",
		enableHeartbeatCheck:   defaultServerHeartbeatCheck,
		heartbeatCheckInterval: defaultServerHeartbeatCheckInterval * time.Second,
		handshakeTimeout:       defaultServerHandshakeTimeout * time.Second,
	}
}

// WithServerListenAddr 设置监听地址
func WithServerListenAddr(addr string) ServerOption {
	return func(o *serverOptions) { o.addr = addr }
}

// WithServerMaxConnNum 设置连接的最大连接数
func WithServerMaxConnNum(maxConnNum int) ServerOption {
	return func(o *serverOptions) { o.maxConnNum = maxConnNum }
}

// WithServerPath 设置Websocket的连接路径
func WithServerPath(path string) ServerOption {
	return func(o *serverOptions) { o.path = path }
}

// WithServerCredentials 设置证书和秘钥
func WithServerCredentials(certFile, keyFile string) ServerOption {
	return func(o *serverOptions) { o.keyFile, o.certFile = keyFile, certFile }
}

// WithServerCheckOrigin 设置Websocket跨域检测函数
func WithServerCheckOrigin(checkOrigin CheckOriginFunc) ServerOption {
	return func(o *serverOptions) { o.checkOrigin = checkOrigin }
}

// WithServerEnableHeartbeatCheck 是否启用心跳检测
func WithServerEnableHeartbeatCheck(enable bool) ServerOption {
	return func(o *serverOptions) { o.enableHeartbeatCheck = enable }
}

// WithServerHeartbeatCheckInterval 设置心跳检测间隔时间
func WithServerHeartbeatCheckInterval(heartbeatCheckInterval time.Duration) ServerOption {
	return func(o *serverOptions) { o.heartbeatCheckInterval = heartbeatCheckInterval }
}

// WithServerHandshakeTimeout 设置握手超时时间
func WithServerHandshakeTimeout(handshakeTimeout time.Duration) ServerOption {
	return func(o *serverOptions) { o.handshakeTimeout = handshakeTimeout }
}

// WithServerMaxMsgLen 设置最大包体长度
func WithServerMaxMsgLen(maxMsgLen int) ServerOption {
	return func(o *serverOptions) { o.maxMsgLen = maxMsgLen }
}
