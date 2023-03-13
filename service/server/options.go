package server

import "github.com/tieing/lemon/discovery"

type Options struct {
	NatsAddr string // nats://127.0.0.1:4222,nats://127.0.0.1:4223
	User     string // 用户名
	Pwd      string // 密码

	ConsulAddr string
	Server     discovery.ServiceInstance
}

type Option func(o *Options)

func WithClusterName(name string) Option {
	return func(o *Options) {

	}
}
