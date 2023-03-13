package ws

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/tieing/lemon/network/inet"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type serverConn struct {
	rw                sync.RWMutex    // 锁
	id                string          // 连接ID
	uid               int             // 用户ID
	state             int32           // 连接状态
	conn              *websocket.Conn // WS源连接
	server            *server         // server
	chWrite           chan chWrite    // 写入队列
	done              chan struct{}   // 写入完成信号
	lastHeartbeatTime int64           // 上次心跳时间
	sessionCache      sync.Map        // 缓存数据
	cacheStorage      inet.MetaStorage
}

// ID 获取连接ID
func (c *serverConn) ID() string {
	return c.id
}

// UID 获取用户ID string
func (c *serverConn) UID() int {
	return c.uid
}

// Bind 绑定用户ID string
func (c *serverConn) Bind(uid int) {
	c.uid = uid
	if c.cacheStorage != nil {
		data, err := c.cacheStorage.Select(c.uid)
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Err(err).Send()
			return
		}
		if len(data) > 0 {
			for s, a := range data {
				c.sessionCache.Store(s, a)
			}
		}
	}
}

// Unbind 解绑用户ID
func (c *serverConn) Unbind() {

	c.sessionCache = sync.Map{}
	if c.cacheStorage != nil {
		_ = c.cacheStorage.Delete(c.uid)
	}
}

// Send 打包并发送消息（同步）
func (c *serverConn) Send(data []byte, msgType ...int) (err error) {
	if err = c.checkState(); err != nil {
		return
	}

	if len(msgType) == 0 {
		msgType = append(msgType, BinaryMessage)
	}
	return c.conn.WriteMessage(msgType[0], data)
}

// Push 打包并发送消息（异步）
func (c *serverConn) Push(data []byte, msgType ...int) (err error) {
	if err = c.checkState(); err != nil {
		return
	}

	if len(msgType) == 0 {
		msgType = append(msgType, BinaryMessage)
	}
	c.chWrite <- chWrite{typ: dataPacket, msg: data, msgType: msgType[0]}
	return nil
}

// PushAll 打包并全网关发送(异步)
func (c *serverConn) PushAll(data []byte, msgType ...int) error {
	if len(msgType) == 0 {
		msgType = append(msgType, BinaryMessage)
	}
	c.server.connMgr.Each(func(conn inet.Conn) bool {
		c.chWrite <- chWrite{typ: dataPacket, msg: data, msgType: msgType[0]}
		return true
	})
	return nil
}

// State 获取连接状态
func (c *serverConn) State() inet.ConnState {
	return inet.ConnState(atomic.LoadInt32(&c.state))
}

// Close 关闭连接
func (c *serverConn) Close(isForce ...bool) error {
	if len(isForce) > 0 && isForce[0] {
		return c.forceClose()
	} else {
		return c.graceClose()
	}
}

// LocalIP 获取本地IP
func (c *serverConn) LocalIP() (string, error) {
	addr, err := c.LocalAddr()
	if err != nil {
		return "", err
	}

	return addr.String(), nil
	//return xnet.ExtractIP(addr)
}

// LocalAddr 获取本地地址
func (c *serverConn) LocalAddr() (net.Addr, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err := c.checkState(); err != nil {
		return nil, err
	}

	return c.conn.LocalAddr(), nil
}

// RemoteIP 获取远端IP
func (c *serverConn) RemoteIP() (string, error) {
	addr, err := c.RemoteAddr()
	if err != nil {
		return "", err
	}

	return addr.String(), nil
	//return xnet.ExtractIP(addr)
}

// RemoteAddr 获取远端地址
func (c *serverConn) RemoteAddr() (net.Addr, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err := c.checkState(); err != nil {
		return nil, err
	}

	return c.conn.RemoteAddr(), nil
}

func (c *serverConn) GetMateData(k string) any {
	v, _ := c.sessionCache.Load(k)
	return v
}
func (c *serverConn) DelMateData(k string) {
	c.sessionCache.Delete(k)
	if c.cacheStorage == nil {
		return
	}

	data := c.GetMetaDataAll()
	if len(data) == 0 { // 删除缓存记录
		_ = c.cacheStorage.Delete(c.uid)
	} else {
		_ = c.cacheStorage.Storage(c.uid, data)
	}

}

func (c *serverConn) SetMetaData(k string, v any) {

	c.sessionCache.Store(k, v)
	if c.cacheStorage != nil {
		_ = c.cacheStorage.Storage(c.uid, c.GetMetaDataAll())
	}
}

func (c *serverConn) GetMetaDataAll() map[string]any {
	data := map[string]any{}
	c.sessionCache.Range(func(key, value any) bool {
		data[key.(string)] = value
		return true
	})
	return data
}

// 初始化连接
func (c *serverConn) init(conn *websocket.Conn, cid string, svr *server) {
	c.id = cid
	c.conn = conn
	c.server = svr
	c.chWrite = make(chan chWrite, 64)
	c.done = make(chan struct{})
	c.cacheStorage = c.server.mateStorage
	c.sessionCache = sync.Map{}
	c.lastHeartbeatTime = time.Now().Unix()
	atomic.StoreInt32(&c.state, int32(inet.ConnOpened))

	if c.server.connectHandler != nil {
		c.server.connectHandler(c)
	}

	go c.read()

	go c.write()
}

// 检测连接状态
func (c *serverConn) checkState() error {
	switch inet.ConnState(atomic.LoadInt32(&c.state)) {
	case inet.ConnHanged:
		return inet.ErrConnectionHanged
	case inet.ConnClosed:
		return inet.ErrConnectionClosed
	}

	return nil
}

// 读取消息
func (c *serverConn) read() {

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			_ = c.forceClose()
			return
		}

		if len(data) > c.server.opts.maxMsgLen {
			log.Warn().Msg("the msg size too large, has been ignored")
			continue
		}

		atomic.StoreInt64(&c.lastHeartbeatTime, time.Now().Unix())

		switch c.State() {
		case inet.ConnHanged:
			continue
		case inet.ConnClosed:
			return
		}

		// ignore heartbeat packet
		if len(data) == 1 {
			c.Push([]byte("0")) // 心跳回复
			continue
		}

		c.server.receiveHandler(c, data)
	}
}

// 优雅关闭
func (c *serverConn) graceClose() (err error) {
	c.rw.Lock()

	if err = c.checkState(); err != nil {
		c.rw.Unlock()
		return
	}

	atomic.StoreInt32(&c.state, int32(inet.ConnHanged))
	c.chWrite <- chWrite{typ: closeSig}
	c.rw.Unlock()

	<-c.done

	c.rw.Lock()
	atomic.StoreInt32(&c.state, int32(inet.ConnClosed))
	close(c.chWrite)
	close(c.done)
	err = c.conn.Close()

	c.conn = nil

	c.server.connMgr.recycle(c)
	c.rw.Unlock()

	if c.server.disconnectHandler != nil {
		c.server.disconnectHandler(c)
	}

	return
}

// 强制关闭
func (c *serverConn) forceClose() (err error) {
	c.rw.Lock()

	if err = c.checkState(); err != nil {
		c.rw.Unlock()
		return
	}

	atomic.StoreInt32(&c.state, int32(inet.ConnClosed))
	close(c.chWrite)
	close(c.done)
	err = c.conn.Close()
	c.conn = nil
	c.server.connMgr.recycle(c)
	c.rw.Unlock()

	if c.server.disconnectHandler != nil {
		c.server.disconnectHandler(c)
	}

	return
}

// 写入消息
func (c *serverConn) write() {
	var ticker *time.Ticker
	if c.server.opts.enableHeartbeatCheck {
		ticker = time.NewTicker(c.server.opts.heartbeatCheckInterval)
		defer ticker.Stop()
	} else {
		ticker = &time.Ticker{C: make(chan time.Time, 1)}
	}

	for {
		select {
		case write, ok := <-c.chWrite:
			if !ok {
				return
			}

			if write.typ == closeSig {
				c.done <- struct{}{}
				return
			}

			if err := c.doWrite(&write); err != nil {
				log.Error().Msgf("write message error: %v", err)
				_ = c.forceClose()
				return
			}
		case <-ticker.C:
			deadline := time.Now().Add(-2 * c.server.opts.heartbeatCheckInterval).Unix()
			if atomic.LoadInt64(&c.lastHeartbeatTime) < deadline {
				log.Debug().Msgf("connection heartbeat timeout")
				_ = c.Close(true)
				return
			}
		}
	}
}

func (c *serverConn) doWrite(write *chWrite) error {
	if atomic.LoadInt32(&c.state) == int32(inet.ConnClosed) {
		return nil
	}

	return c.conn.WriteMessage(write.msgType, write.msg)
}
