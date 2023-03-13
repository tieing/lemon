package discovery

import (
	"context"
)

// State 集群实例状态
type State string

const (
	Work State = "work" // 工作（节点正常工作，可以分配更多玩家到该节点）
	Busy State = "busy" // 繁忙（节点资源紧张，不建议分配更多玩家到该节点上）
	Hang State = "hang" // 挂起（节点即将关闭，正处于资源回收中）
	Done State = "done" // 关闭（节点已经关闭，无法正常访问该节点）
)

type Registry interface {
	// KV 数据存储
	KV
	// Register 注册服务实例
	Register(ins *ServiceInstance) error
	// Deregister 解注册服务实例
	Deregister(ctx context.Context, ins *ServiceInstance) error
	// Watch 监听相同服务名的服务实例变化
	Watch(ctx context.Context, serviceName string, cb func([]*ServiceInstance)) (Watcher, error)
	// Services 获取服务实例列表
	Services(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
}

type Watcher interface {
	// ServiceInstances 返回服务实例列表
	Services() []*ServiceInstance
	// Stop 停止监听
	Stop()
}

/*
服务ID:"xxxx_01","xxxx_02","xxxx_02"
服务名称:"game","agent","gate","webapi",""chat,"gm","log" // 一级分类
服务类别:"crash","dice","gate","webapi","gm","log" // 二级分类
服务别名:"crash01服","dice01服",
服务状态:""
服务端口:""(地址(不通过ws网关调用的服务),订阅/发布的字符串(通过网关调用的服务))

------------
游戏多服配置:
服务名称:'游戏类'固定为"game"
服务类别:当前游戏的类别,相同游戏服务类别一样
服务ID: 使用 服务类别+"_01" "_02" 来区分

*/

type ServiceInstance struct {
	// 服务实体ID，每个服务实体ID唯一
	ID string `json:"id"`
	// 服务实体名(如:login,game,DB,lobby,logger)
	Name string `json:"name"`
	// 服务实体类型
	Kind string `json:"kind"`
	// 服务实体别名
	Alias string `json:"alias"`
	// 服务实例状态
	State State `json:"state"`
	// 权重 (多个相同服务才生效)
	Weight int `json:"weighted"`
	// 服务器实体暴露端口
	Endpoint string `json:"endpoint                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             "`
}
