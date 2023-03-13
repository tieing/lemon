package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	registry "github.com/tieing/lemon/discovery"
	"strconv"
	"time"
)

const (
	checkIDFormat     = "service:%s"
	checkUpdateOutput = "passed"
	metaFieldKind     = "kind"
	metaFieldAlias    = "alias"
	metaFieldState    = "state"
	metaFieldWeight   = "weight"
)

type registrar struct {
	ctx      context.Context
	cancel   context.CancelFunc
	registry *Registry
}

func newRegistrar(ctx context.Context, registry *Registry) *registrar {
	r := &registrar{}
	r.ctx, r.cancel = context.WithCancel(registry.ctx)
	r.registry = registry

	return r
}

// 注册服务
func (r *registrar) register(ctx context.Context, ins *registry.ServiceInstance) error {
	registration := &api.AgentServiceRegistration{
		ID:      ins.ID,
		Name:    ins.Name,
		Meta:    make(map[string]string),
		Address: ins.Endpoint,
	}

	registration.Meta[metaFieldKind] = ins.Kind
	registration.Meta[metaFieldAlias] = ins.Alias
	registration.Meta[metaFieldState] = string(ins.State)
	registration.Meta[metaFieldWeight] = strconv.Itoa(ins.Weight)

	// 声明心跳检测方式
	if r.registry.opts.EnableHealthCheck {
		registration.Check = &api.AgentServiceCheck{
			CheckID:                        fmt.Sprintf(checkIDFormat, ins.ID),
			TTL:                            fmt.Sprintf("%ds", r.registry.opts.HealthCheckInterval),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", r.registry.opts.CheckFailedDeregister),
		}
	}

	if err := r.registry.client.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	if r.registry.opts.EnableHealthCheck {
		// 启动健康心跳
		go r.heartbeat(ctx, ins.ID)
	}

	return nil
}

// 解注册服务
func (r *registrar) deregister(ctx context.Context, ins *registry.ServiceInstance) error {
	r.cancel()

	r.registry.registrars.Delete(ins.ID)

	return r.registry.client.Agent().ServiceDeregister(ins.ID)
}

// 心跳
func (r *registrar) heartbeat(ctx context.Context, insID string) {
	checkID := fmt.Sprintf(checkIDFormat, insID)

	ticker := time.NewTicker(time.Duration(r.registry.opts.HealthCheckInterval) * time.Second / 2)
	defer ticker.Stop()
	var err error
	for {
		select {
		case <-ticker.C:
			if err = r.registry.client.Agent().UpdateTTL(checkID, checkUpdateOutput, api.HealthPassing); err != nil {
				log.Err(err).Msg("update heartbeat ttl failed!!")
			}
		case <-ctx.Done():
			return
		}
	}
}
