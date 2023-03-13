package consul

import (
	"context"
	"github.com/tieing/lemon/discovery"
	"sync"
	"sync/atomic"
	"time"
)

type watcher struct {
	err              error
	ctx              context.Context
	cancel           context.CancelFunc
	discovery        *Registry
	serviceName      string
	serviceInstances *atomic.Value
	serviceWaitIndex uint64
	idx              int64
	watcherNum       int32
	rw               sync.RWMutex
	cb               func(services []*discovery.ServiceInstance)
}

func newWatcher(ctx context.Context, serviceName string, discovery *Registry) (*watcher, error) {
	services, index, err := discovery.services(ctx, serviceName, 0, true)
	if err != nil {
		return nil, err
	}

	wm := &watcher{}
	wm.ctx, wm.cancel = context.WithCancel(discovery.ctx)
	wm.discovery = discovery
	wm.serviceName = serviceName
	wm.serviceInstances = &atomic.Value{}
	wm.serviceWaitIndex = index
	wm.serviceInstances.Store(services)
	wm.watcherNum = 1

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-wm.ctx.Done():
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
				services, index, err = wm.discovery.services(ctx, wm.serviceName, wm.serviceWaitIndex, true)
				cancel()
				if err != nil {
					time.Sleep(time.Second)
					continue
				}

				if index != wm.serviceWaitIndex {
					wm.serviceWaitIndex = index
					wm.serviceInstances.Store(services)
					wm.broadcast()
				}
			}
		}
	}()
	return wm, nil
}

func (wm *watcher) broadcast() {
	if wm.cb != nil {
		wm.rw.RLock()
		defer wm.rw.RUnlock()
		services := wm.Services()
		wm.cb(services)
	}
}

func (wm *watcher) Services() []*discovery.ServiceInstance {
	return wm.serviceInstances.Load().([]*discovery.ServiceInstance)
}

func (wm *watcher) Stop() {
	if atomic.AddInt32(&wm.watcherNum, -1) == 0 {
		wm.cancel()
	}
}
