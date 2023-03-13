package tcp

import (
	"github.com/rs/zerolog/log"
	"github.com/tieing/lemon/network/inet"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type serverConn struct {
	rw                sync.RWMutex  // 锁
	id                string        // 连接ID
	uid               int           // 用户ID
	state             int32         // 连接状态
	conn              net.Conn      // TCP源连接
	server            *server       // server
	chWrite           chan chWrite  // 写入队列
	lastHeartbeatTime int64         // 上次心跳时间
	done              chan struct{} // 写入完成信号
	sessionCache      sync.Map      // 缓存数据
	cacheStorage      inet.MetaStorage
}

var _ inet.Conn = &serverConn{}

// ID 获取连接ID
func (c *serverConn) ID() string {
	return c.id
}

// UID 获取用户ID
func (c *serverConn) UID() int {
	return c.uid
}

// Bind 绑定用户ID
func (c *serverConn) Bind(uid int) {
	c.uid = uid
}

// Unbind 解绑用户ID
func (c *serverConn) Unbind() {
	c.uid = 0
}

// Send 发送消息（同步）
func (c *serverConn) Send(msg []byte, msgType ...int) (err error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err = c.checkState(); err != nil {
		return
	}

	_, err = c.conn.Write(msg)

	return
}

// Push 发送消息（异步）
func (c *serverConn) Push(msg []byte, msgType ...int) (err error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err = c.checkState(); err != nil {
		return
	}

	c.chWrite <- chWrite{typ: dataPacket, msg: msg}

	return
}

// PushAll 全网关发送(异步)
func (c *serverConn) PushAll(data []byte, msgType ...int) error {

	c.server.connMgr.Each(func(conn inet.Conn) bool {
		c.chWrite <- chWrite{typ: dataPacket, msg: data}
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

// 初始化连接
func (c *serverConn) init(conn net.Conn, cid string, svr *server) {
	c.id = cid
	c.uid = 0
	c.conn = conn
	c.server = svr
	c.chWrite = make(chan chWrite, 1024)
	c.done = make(chan struct{})
	c.lastHeartbeatTime = time.Now().Unix()
	atomic.StoreInt32(&c.state, int32(inet.ConnOpened))

	if c.server.connectHandler != nil {
		c.server.connectHandler(c)
	}

	go c.read()

	go c.write()
}

// 读取消息
func (c *serverConn) read() {
	for {
		msg, err := readMsgFromConn(c.conn, c.server.opts.maxMsgLen)
		if err != nil {
			if err == errMsgSizeTooLarge {
				log.Warn().Msg("the msg size too large, has been ignored")
				continue
			}
			_ = c.forceClose()
			break
		}

		atomic.StoreInt64(&c.lastHeartbeatTime, time.Now().Unix())

		switch c.State() {
		case inet.ConnHanged:
			continue
		case inet.ConnClosed:
			return
		}

		// ignore heartbeat packet
		if len(msg) == 0 {
			continue
		}

		if c.server.receiveHandler != nil {
			c.server.receiveHandler(c, msg)
		}
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

			buf, err := pack(write.msg)
			if err != nil {
				log.Error().Msgf("packet message error: %v", err)
				continue
			}

			if err = c.doWrite(buf); err != nil {
				log.Error().Msgf("write message error: %v", err)
			}
		case <-ticker.C:
			deadline := time.Now().Add(-2 * c.server.opts.heartbeatCheckInterval).Unix()
			if atomic.LoadInt64(&c.lastHeartbeatTime) < deadline {
				log.Debug().Msg("connection heartbeat timeout")
				_ = c.Close(true)
				return
			}
		}
	}
}

func (c *serverConn) doWrite(buf []byte) (err error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if atomic.LoadInt32(&c.state) == int32(inet.ConnClosed) {
		return
	}

	_, err = c.conn.Write(buf)

	return
}
