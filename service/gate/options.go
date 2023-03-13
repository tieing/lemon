package gate

import "github.com/tieing/lemon/discovery"

type Options struct {
	NatsAddr string // nats://127.0.0.1:4222,nats://127.0.0.1:4223
	User     string // 用户名
	Pwd      string // 密码

	ListenAddr string // 监听地址
	MaxMsgLen  int    // 最大消息长度（字节），默认1kb
	MaxConnNum int    // 最大连接数
	CertFile   string // 证书文件
	KeyFile    string // 秘钥文件
	ConsulAddr string
	Gate       discovery.ServiceInstance
	Agent      discovery.ServiceInstance
}

type Option func(o *Options)

func WithClusterName(name string) Option {
	return func(o *Options) {

	}
}
