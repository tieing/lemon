package gate

import (
	"context"
	"fmt"
	"github.com/tieing/lemon/discovery"
	"github.com/tieing/lemon/discovery/consul"
	"github.com/tieing/lemon/network/inet"
	"github.com/tieing/lemon/network/ws"
	"github.com/tieing/lemon/processor"
	"github.com/tieing/lemon/rpc"
	"github.com/tieing/lemon/rpc/rpc_impl"
	"github.com/tieing/lemon/rpc/rpc_impl/rpcmsg"
	"github.com/tieing/lemon/selector"
	"github.com/tieing/lemon/service/gate/gatemsg"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

const (
	bindRemoteServiceKey = "__remote_server"
)

type Gate struct {
	ctx        context.Context
	cancel     context.CancelFunc
	serverID   string
	opt        *Options
	ws         inet.SocketServer
	rpc        rpc.RPC
	registry   discovery.Registry
	netMsgPool sync.Pool

	mx        sync.Mutex
	services  sync.Map //map[string][]*discovery.ServiceInstance // 服务
	watcher   map[string]discovery.Watcher
	selectors sync.Map //map[string]*weighted.SW
	//netMsgQue *queue.Queue

	processor processor.Processor
}

func NewGate(ctx context.Context, serverID string, process processor.Processor, opt *Options) *Gate {
	ctx, cancel := context.WithCancel(ctx)
	gate := &Gate{
		ctx:        ctx,
		cancel:     cancel,
		serverID:   serverID,
		opt:        opt,
		netMsgPool: sync.Pool{New: func() any { return &gatemsg.NetworkCMDMessage{} }},
		watcher:    map[string]discovery.Watcher{},
		//netMsgQue:  queue.New(),
		processor: process,
	}

	gate.ws = ws.NewServer(
		serverID,
		ws.WithServerListenAddr(opt.ListenAddr),
		ws.WithServerMaxMsgLen(opt.MaxMsgLen),
		ws.WithServerEnableHeartbeatCheck(true),
		ws.WithServerMaxConnNum(opt.MaxConnNum),
	)
	gate.ws.OnStart(func() {
		fmt.Printf("server on start. time:%s", time.Now().Format("2006-01-02 15:04:05"))
	})
	gate.ws.OnStop(func() {
		fmt.Printf("server on stop. time:%s", time.Now().Format("2006-01-02 15:04:05"))
	})
	gate.ws.OnConnect(func(conn inet.Conn) {
		println("连接打开!!!! CID: " + conn.ID())
	})
	gate.ws.OnDisconnect(gate.doNetDisconnect)
	gate.ws.OnReceive(gate.doNetMessage)

	// 注册内部方法
	gate.registerDefaultHandler()
	return gate
}

func (p *Gate) Run(watch ...string) error {
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

	err = p.registry.Register(&p.opt.Gate)
	if err != nil {
		return err
	}

	err = p.registry.Register(&p.opt.Agent)
	if err != nil {
		return err
	}
	err = p.WatchService(watch)
	if err != nil {
		zap.S().Fatal(err, "服务发现启动失败! 进程退出!")
	}

	return p.ws.Start()
}

func (p *Gate) Close() {
	p.ws.Stop()
	p.rpc.Close()
	p.cancel()
}

func (p *Gate) doNetMessage(conn inet.Conn, buf []byte) {
	msg := gatemsg.NetworkCMDMessage{}
	err := proto.Unmarshal(buf, &msg)
	if err != nil {
		fmt.Printf("解码网络数据流失败!, err : %s", err.Error())
		return
	}
	var serverID string
	if v := conn.GetMateData(bindRemoteServiceKey); v != nil {
		serverID = v.(string)
	} else {
		c, ok := p.selectors.Load(msg.Service)
		if !ok {
			fmt.Printf("查找服务失败!, server Name : %s", msg.Service)
			return
		}
		item := c.(*selector.SW).Next()
		if item == nil {
			fmt.Printf("没有可用的服务 !, server Name : %s", msg.Service)
			return
		}
		serverID = item.(*discovery.ServiceInstance).ID
		conn.SetMetaData(bindRemoteServiceKey, serverID)
	}

	// 解析 msg 获得远程 服务类型
	// 根据远程服务类型进行负载均衡选择出服务ID
	// 获得服务ID后将数据转发到对应服务

	p.rpc.Push(&rpcmsg.RPCMessage{
		Type: rpcmsg.Type_C2S,
		MID:  msg.Mid,
		SID:  p.serverID,
		CID:  conn.ID(),
		UID:  int32(conn.UID()),
		BUF:  msg.Buf,
		RID:  serverID,
		SEQ:  0,
	})
}

func (p *Gate) doNetDisconnect(conn inet.Conn) {
	println("连接断开! CID:" + conn.ID())
	if v := conn.GetMateData(bindRemoteServiceKey); v != nil {
		serverID := v.(string)
		p.rpc.Push(&rpcmsg.RPCMessage{
			Type: rpcmsg.Type_Notify,
			HID:  "__OnSessionDisconnect",
			SID:  p.serverID,
			CID:  conn.ID(),
			UID:  int32(conn.UID()),
			RID:  serverID,
		})
	}
}

// 注册本地函数
// session操作函数
// 控制函数

func (r *Gate) WatchService(serverNames []string) error {
	for _, name := range serverNames {
		watch, err := r.registry.Watch(r.ctx, name, func(instance []*discovery.ServiceInstance) {
			r.mx.Lock()
			// 服务改变,先清理原来的服务,然后存入最新的服务
			serviceGroup := map[string][]*discovery.ServiceInstance{}
			for _, serviceInstance := range instance {
				serviceGroup[serviceInstance.Kind] = append(serviceGroup[serviceInstance.Kind], serviceInstance)
			}
			for s, instList := range serviceGroup {
				r.services.Store(s, instList)
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
			r.mx.Unlock()
		})
		if err != nil {
			return err
		}

		r.watcher[name] = watch
	}
	return nil
}
