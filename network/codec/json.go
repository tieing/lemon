package codec

import "encoding/json"

type JsonCodec struct{}

func (*JsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (*JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
