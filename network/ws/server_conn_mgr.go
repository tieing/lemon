package ws

import (
	"fmt"
	"github.com/tieing/lemon/network/inet"
	"math"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type connMgr struct {
	count    int32     // 连接数量
	id       int32     // 连接ID
	serverID string    // 连接ID前缀(服务ID,避免多网关时连接id重复)
	pool     sync.Pool // 连接池
	conns    sync.Map  //map[connID]*serverConn // 连接集合
}

func newConnMgr(server *server) *connMgr {
	return &connMgr{
		conns:    sync.Map{},
		serverID: server.serverID,
		pool:     sync.Pool{New: func() interface{} { return &serverConn{} }},
	}
}

// 关闭连接
func (cm *connMgr) close() {
	cm.conns.Range(func(_, value any) bool {
		value.(*serverConn).Close(false)
		return true
	})
	atomic.StoreInt32(&cm.count, 0)
}

// 分配连接
func (cm *connMgr) allocate(server *server, c *websocket.Conn) {
	atomic.AddInt32(&cm.count, 1)
	cid := fmt.Sprintf("%s_%d", cm.serverID, atomic.AddInt32(&cm.id, 1))
	conn := cm.pool.Get().(*serverConn)
	cm.conns.Store(cid, conn)

	conn.init(c, cid, server)
}

// 回收连接
func (cm *connMgr) recycle(conn *serverConn) {
	cm.conns.LoadAndDelete(conn.ID())
	cm.pool.Put(conn)
	if atomic.AddInt32(&cm.count, -1) > math.MaxInt32-1000 {
		atomic.StoreInt32(&cm.id, 1)
	}
}

// Each 遍历连接
func (cm *connMgr) Each(fn func(conn inet.Conn) bool) {
	cm.conns.Range(func(_, value any) bool {
		return fn(value.(inet.Conn))
	})
}

func (cm *connMgr) GetConn(id string) inet.Conn {
	if v, ok := cm.conns.Load(id); ok {
		return v.(inet.Conn)
	}
	return nil
}
