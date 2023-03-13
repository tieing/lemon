package packet

import (
	"encoding/binary"
	"errors"
	"github.com/tieing/lemon/network/inet"
)

var (
	ErrMessageIsNil  = errors.New("the message is nil")
	ErrRouteOverflow = errors.New("the message route overflow")
)

type Packet struct{}

var globalPacker = NewPacker()

func NewPacker() inet.Packer {
	return &Packet{}
}

// Pack 打包消息
func (p *Packet) Pack(mid int32, data []byte) ([]byte, error) {

	if mid > int32(1<<(8*4-1)-1) || mid < int32(-1<<(8*4-1)) {
		return nil, ErrRouteOverflow
	}

	buf := make([]byte, len(data)+4)
	binary.LittleEndian.PutUint32(buf[:4], uint32(mid))
	copy(buf[4:], data)
	return buf, nil
}

// Unpack 解包消息
func (p *Packet) Unpack(buf []byte) (mid int32, data []byte, err error) {
	if len(buf) < 4 {
		return 0, nil, ErrMessageIsNil
	}
	msgId := int32(binary.LittleEndian.Uint32(buf[:4]))

	return msgId, buf[4:], nil
}

// 全局打包消息

// Pack 打包消息
func Pack(mid int32, data []byte) ([]byte, error) {
	return globalPacker.Pack(mid, data)
}

// Unpack 解包消息
func Unpack(buf []byte) (mid int32, data []byte, err error) {
	return globalPacker.Unpack(buf)
}
