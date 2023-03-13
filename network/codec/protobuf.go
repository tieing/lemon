package codec

import "google.golang.org/protobuf/proto"

var CodecProtobuf = NewProtobufCodec()

type ProtobufCodec struct{}

func NewProtobufCodec() *ProtobufCodec {
	return new(ProtobufCodec)
}
func (ProtobufCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (ProtobufCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}
