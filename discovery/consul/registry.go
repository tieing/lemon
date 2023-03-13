package consul

import (
	"context"
	"github.com/tieing/lemon/discovery"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
)

type Registry struct {
	ctx        context.Context
	cancel     context.CancelFunc
	opts       *Options
	client     *api.Client
	watchers   sync.Map
	registrars sync.Map
}

func NewRegistry(ctx context.Context, opts ...Option) (*Registry, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	r := &Registry{}
	r.opts = o
	r.ctx, r.cancel = context.WithCancel(ctx)

	config := api.DefaultConfig()
	config.Address = o.Addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	r.client = client
	return r, nil
}

// Register 注册服务实例
func (r *Registry) Register(ins *discovery.ServiceInstance) error {

	v, ok := r.registrars.Load(ins.ID)
	if ok {
		return v.(*registrar).register(r.ctx, ins)
	}

	reg := newRegistrar(r.ctx, r)
	if err := reg.register(r.ctx, ins); err != nil {
		return err
	}
	r.registrars.Store(ins.ID, reg)

	return nil
}

// Deregister 解注册服务实例
func (r *Registry) Deregister(ctx context.Context, ins *discovery.ServiceInstance) error {
	v, ok := r.registrars.Load(ins.ID)
	if ok {
		return v.(*registrar).deregister(ctx, ins)
	}

	return r.client.Agent().ServiceDeregister(ins.ID)
}

// Services 获取服务实例列表
func (r *Registry) Services(ctx context.Context, serviceName string) ([]*discovery.ServiceInstance, error) {
	v, ok := r.watchers.Load(serviceName)
	if ok {
		return v.(*watcher).Services(), nil
	} else {
		services, _, err := r.services(ctx, serviceName, 0, true)
		return services, err
	}
}

// Watch 监听服务
func (r *Registry) Watch(ctx context.Context, serviceName string, cb func([]*discovery.ServiceInstance)) (discovery.Watcher, error) {
	v, ok := r.watchers.Load(serviceName)
	if ok {
		w := v.(*watcher)
		atomic.AddInt32(&w.watcherNum, 1)
		return w, nil
	}

	w, err := newWatcher(ctx, serviceName, r)
	if err != nil {
		return nil, err
	}

	w.cb = cb
	r.watchers.Store(serviceName, w)

	if w.cb != nil {
		w.cb(w.Services())
	}
	return w, nil
}

// 获取服务实体列表
func (r *Registry) services(ctx context.Context, serviceName string, waitIndex uint64, passingOnly bool) ([]*discovery.ServiceInstance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: waitIndex,
		WaitTime:  60 * time.Second,
	}
	opts.WithContext(ctx)

	entries, meta, err := r.client.Health().Service(serviceName, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}

	services := make([]*discovery.ServiceInstance, 0, len(entries))
	for _, entry := range entries {

		if entry.Service.ID == "" {
			continue
		}
		ins := &discovery.ServiceInstance{
			ID:       entry.Service.ID,
			Name:     entry.Service.Service,
			Endpoint: entry.Service.Address,
		}
		for k, v := range entry.Service.Meta {
			switch k {
			case metaFieldKind:
				ins.Kind = v
			case metaFieldAlias:
				ins.Alias = v
			case metaFieldState:
				ins.State = discovery.State(v)
			case metaFieldWeight:
				ins.Weight, _ = strconv.Atoi(v)
			default:
				continue
			}
		}
		services = append(services, ins)
	}
	return services, meta.LastIndex, nil
}
