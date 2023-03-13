package server

import (
	"context"
	"github.com/tieing/lemon/discovery"
	"github.com/tieing/lemon/discovery/consul"
	"github.com/tieing/lemon/processor"
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl"
	"github.com/tieing/lemon/selector"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"sync"
)

type Server struct {
	ctx       context.Context
	cancel    context.CancelFunc
	serverID  string
	opt       *Options
	rpc       rpc.RPC
	registry  discovery.Registry
	mx        sync.Mutex
	services  sync.Map //map[instanceID]*discovery.ServiceInstance // 服务
	watcher   map[string]discovery.Watcher
	selectors sync.Map //map[string]*weighted.SW
	processor processor.Processor

	// method
	onSessionDisconnect func(gateID, CID string, uid int)
}

func NewServer(ctx context.Context, serverID string, pro processor.Processor, opt *Options) *Server {
	ctx, cancel := context.WithCancel(ctx)
	gate := &Server{
		ctx:       ctx,
		cancel:    cancel,
		serverID:  serverID,
		opt:       opt,
		watcher:   map[string]discovery.Watcher{},
		processor: pro,
	}

	// 注册内部方法
	gate.registerDefaultHandler()
	return gate
}

func (p *Server) Run(watch ...string) error {
	conn, err := nats.Connect(p.opt.NatsAddr, nats.UserInfo(p.opt.User, p.opt.Pwd))
	if err != nil {
		return err
	}
	p.rpc = rpc_impl.NewRPC(p.ctx, p.serverID, p.processor, conn)
	err = p.rpc.Run()
	if err != nil {
		return err
	}

	p.registry, err = consul.NewRegistry(p.ctx, consul.WithAddr(p.opt.ConsulAddr))
	if err != nil {
		return err
	}
	err = p.registry.Register(&p.opt.Server)
	if err != nil {
		return err
	}

	err = p.WatchService(watch)
	if err != nil {
		zap.S().Fatal(err, "服务发现启动失败! 进程退出!")
	}

	return nil
}

func (p *Server) GetServerByID(serverID string) rpc.Server {
	if _, ok := p.services.Load(serverID); ok {
		return rpc_impl.NewServer(p.rpc, serverID, 0)
	}
	return nil
}

// 多个服务 返回随机找到的第一个
func (p *Server) GetServerByKind(kind string) (s rpc.Server) {
	p.services.Range(func(_, value any) bool {
		if value.(*discovery.ServiceInstance).Kind == kind {
			s = rpc_impl.NewServer(p.rpc, value.(*discovery.ServiceInstance).ID, 0)
			return false
		}
		return true
	})
	return nil
}

func (p *Server) SetOnSessionDisconnect(f func(gateID, CID string, uid int)) {
	p.onSessionDisconnect = f
}

func (p *Server) Close() {
	p.rpc.Close()
	p.cancel()
}

// 注册本地函数
// session操作函数
// 控制函数

func (r *Server) WatchService(serverNames []string) error {
	for _, name := range serverNames {
		watch, err := r.registry.Watch(r.ctx, name, func(instance []*discovery.ServiceInstance) {
			// 服务改变,先清理原来的服务,然后存入最新的服务
			serviceGroup := map[string][]*discovery.ServiceInstance{}
			services := sync.Map{}
			for _, serviceInstance := range instance {
				services.Store(serviceInstance.ID, serviceInstance)
				serviceGroup[serviceInstance.Kind] = append(serviceGroup[serviceInstance.Kind], serviceInstance)
			}

			for s, instList := range serviceGroup {
				var lbSw *selector.SW
				if sw, ok := r.selectors.Load(s); ok {
					lbSw = sw.(*selector.SW)
					lbSw.Reset()
				} else {
					lbSw = &selector.SW{}
				}
				for _, d := range instList {
					lbSw.Add(d, d.Weight)
				}
				r.selectors.Store(s, lbSw)
			}
			r.services = services
		})
		if err != nil {
			return err
		}

		r.watcher[name] = watch
	}
	return nil
}
