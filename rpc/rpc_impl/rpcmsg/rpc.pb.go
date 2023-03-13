// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: rpc.proto

package rpcmsg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Type int32

const (
	Type_Invalid  Type = 0
	Type_Request  Type = 1
	Type_Response Type = 2
	Type_Notify   Type = 3
	Type_C2S      Type = 4 // ws客户端到后端服务器的消息
	Type_S2C      Type = 5 // 推送到ws客户端
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "Invalid",
		1: "Request",
		2: "Response",
		3: "Notify",
		4: "C2S",
		5: "S2C",
	}
	Type_value = map[string]int32{
		"Invalid":  0,
		"Request":  1,
		"Response": 2,
		"Notify":   3,
		"C2S":      4,
		"S2C":      5,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_rpc_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{0}
}

// RPC 消息结构
type RPCMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type Type   `protobuf:"varint,1,opt,name=Type,proto3,enum=Type" json:"Type,omitempty"` // [required]
	MID  int32  `protobuf:"varint,2,opt,name=MID,proto3" json:"MID,omitempty"`             // msg ID 消息ID[用于c2s]
	HID  string `protobuf:"bytes,3,opt,name=HID,proto3" json:"HID,omitempty"`              // handlerID 内部消息处理器ID
	BUF  []byte `protobuf:"bytes,4,opt,name=BUF,proto3" json:"BUF,omitempty"`              // buffer 数据
	SID  string `protobuf:"bytes,5,opt,name=SID,proto3" json:"SID,omitempty"`              // sender ID 发送方ID [required]
	CID  string `protobuf:"bytes,6,opt,name=CID,proto3" json:"CID,omitempty"`              // connection ID 客户端连接ID [客户端消息时]
	UID  int32  `protobuf:"varint,7,opt,name=UID,proto3" json:"UID,omitempty"`             // user ID
	SEQ  int32  `protobuf:"varint,8,opt,name=SEQ,proto3" json:"SEQ,omitempty"`             // 请求ID,回执需要
	RID  string `protobuf:"bytes,9,opt,name=RID,proto3" json:"RID,omitempty"`              // 远程服务器ID
	Err  string `protobuf:"bytes,10,opt,name=Err,proto3" json:"Err,omitempty"`             // error
}

func (x *RPCMessage) Reset() {
	*x = RPCMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPCMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPCMessage) ProtoMessage() {}

func (x *RPCMessage) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPCMessage.ProtoReflect.Descriptor instead.
func (*RPCMessage) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *RPCMessage) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_Invalid
}

func (x *RPCMessage) GetMID() int32 {
	if x != nil {
		return x.MID
	}
	return 0
}

func (x *RPCMessage) GetHID() string {
	if x != nil {
		return x.HID
	}
	return ""
}

func (x *RPCMessage) GetBUF() []byte {
	if x != nil {
		return x.BUF
	}
	return nil
}

func (x *RPCMessage) GetSID() string {
	if x != nil {
		return x.SID
	}
	return ""
}

func (x *RPCMessage) GetCID() string {
	if x != nil {
		return x.CID
	}
	return ""
}

func (x *RPCMessage) GetUID() int32 {
	if x != nil {
		return x.UID
	}
	return 0
}

func (x *RPCMessage) GetSEQ() int32 {
	if x != nil {
		return x.SEQ
	}
	return 0
}

func (x *RPCMessage) GetRID() string {
	if x != nil {
		return x.RID
	}
	return ""
}

func (x *RPCMessage) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type PushGroupMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data   []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Member []string `protobuf:"bytes,2,rep,name=member,proto3" json:"member,omitempty"`
}

func (x *PushGroupMsg) Reset() {
	*x = PushGroupMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushGroupMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushGroupMsg) ProtoMessage() {}

func (x *PushGroupMsg) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushGroupMsg.ProtoReflect.Descriptor instead.
func (*PushGroupMsg) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *PushGroupMsg) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PushGroupMsg) GetMember() []string {
	if x != nil {
		return x.Member
	}
	return nil
}

var File_rpc_proto protoreflect.FileDescriptor

var file_rpc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc9, 0x01, 0x0a, 0x0a,
	0x52, 0x50, 0x43, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x05, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x4d, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x48, 0x49, 0x44, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x48, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x42, 0x55, 0x46,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x42, 0x55, 0x46, 0x12, 0x10, 0x0a, 0x03, 0x53,
	0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x53, 0x49, 0x44, 0x12, 0x10, 0x0a,
	0x03, 0x43, 0x49, 0x44, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x43, 0x49, 0x44, 0x12,
	0x10, 0x0a, 0x03, 0x55, 0x49, 0x44, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x55, 0x49,
	0x44, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x45, 0x51, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x53, 0x45, 0x51, 0x12, 0x10, 0x0a, 0x03, 0x52, 0x49, 0x44, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x52, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x72, 0x72, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x45, 0x72, 0x72, 0x22, 0x3a, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x2a, 0x4c, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x49,
	0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x10, 0x03, 0x12,
	0x07, 0x0a, 0x03, 0x43, 0x32, 0x53, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x32, 0x43, 0x10,
	0x05, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x72, 0x70, 0x63, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_proto_rawDescOnce sync.Once
	file_rpc_proto_rawDescData = file_rpc_proto_rawDesc
)

func file_rpc_proto_rawDescGZIP() []byte {
	file_rpc_proto_rawDescOnce.Do(func() {
		file_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_proto_rawDescData)
	})
	return file_rpc_proto_rawDescData
}

var file_rpc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_proto_goTypes = []interface{}{
	(Type)(0),            // 0: Type
	(*RPCMessage)(nil),   // 1: RPCMessage
	(*PushGroupMsg)(nil), // 2: PushGroupMsg
}
var file_rpc_proto_depIdxs = []int32{
	0, // 0: RPCMessage.Type:type_name -> Type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_proto_init() }
func file_rpc_proto_init() {
	if File_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPCMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushGroupMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_proto_depIdxs,
		EnumInfos:         file_rpc_proto_enumTypes,
		MessageInfos:      file_rpc_proto_msgTypes,
	}.Build()
	File_rpc_proto = out.File
	file_rpc_proto_rawDesc = nil
	file_rpc_proto_goTypes = nil
	file_rpc_proto_depIdxs = nil
}