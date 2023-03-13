package inet

type ConnMgr interface {
	Each(fn func(conn Conn) bool)
	GetConn(id string) Conn
}
