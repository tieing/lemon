package consul_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/tieing/lemon/discovery"
	"github.com/tieing/lemon/discovery/consul"
	"github.com/tieing/lemon/tools/random"
	"github.com/tieing/lemon/tools/xnet"
	"net"
	"testing"
	"time"
)

var (
	port = random.RandInt(8000, 10000)
)

const (
	serviceName = "node"
)

var reg, _ = consul.NewRegistry(context.Background())

func server(t *testing.T) {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}

	go func(ls net.Listener) {
		for {
			conn, err := ls.Accept()
			if err != nil {
				t.Error(err)
				return
			}
			var buff []byte
			if _, err = conn.Read(buff); err != nil {
				t.Error(err)
			}
		}
	}(ls)
}

func TestRegistry_Register(t *testing.T) {
	server(t)

	host, err := xnet.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}

	ins := &discovery.ServiceInstance{
		ID:       fmt.Sprintf("%s_%s:%d", serviceName, host, port),
		Name:     serviceName,
		Kind:     "login_service",
		Alias:    "mahjong",
		State:    discovery.Work,
		Endpoint: fmt.Sprintf("grpc://:%d", port),
	}

	if err = reg.Register(ins); err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	ins.State = discovery.Busy
	if err = reg.Register(ins); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300000000 * time.Second)
}
func TestRegistry_Register2(t *testing.T) {
	server(t)
	//serviceName := "node2"
	host, err := xnet.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}

	ins := &discovery.ServiceInstance{
		ID:       fmt.Sprintf("%s_%s:%d", serviceName, host, port),
		Name:     serviceName,
		Kind:     "login_service",
		Alias:    "mahjong",
		State:    discovery.Work,
		Endpoint: fmt.Sprintf("grpc://:%d", port),
	}

	if err = reg.Register(ins); err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	ins.State = discovery.Busy
	if err = reg.Register(ins); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300000000 * time.Second)
}

func TestRegistry_Services(t *testing.T) {
	services, err := reg.Services(context.Background(), serviceName)
	if err != nil {
		t.Fatal(err)
	}

	for _, service := range services {
		t.Logf("%+v", service)
	}
}

func TestRegistry_Watch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	watcher1, err := reg.Watch(ctx, serviceName, nil)
	if err != nil {
		t.Fatal(err)
	}

	watcher2, err := reg.Watch(context.Background(), serviceName, nil)
	if err != nil {
		t.Fatal(err)
	}

	//go func() {
	//	time.Sleep(5 * time.Second)
	//	watcher1.Stop()
	//	time.Sleep(5 * time.Second)
	//	watcher2.Stop()
	//}()

	go func() {
		for {
			services := watcher1.Services()

			fmt.Println("goroutine 1: new event entity")

			for _, service := range services {
				t.Logf("goroutine 1: %+v", service)
			}
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			services := watcher2.Services()

			fmt.Println("goroutine 2: new event entity")

			for _, service := range services {
				t.Logf("goroutine 2: %+v", service)
			}
			time.Sleep(time.Second * 3)
		}
	}()

	time.Sleep(600 * time.Second)
}

func TestRegistry_Nodes(t *testing.T) {
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	h, q, err := client.Health().Service("node", "", true, nil)
	if err != nil {
		panic(err)
	}
	println(h, q)
}
