package inet

type Packer interface {
	// Pack 打包消息
	Pack(mid int32, data []byte) ([]byte, error)
	// Unpack 解包消息
	Unpack(data []byte) (int32, []byte, error)
}
