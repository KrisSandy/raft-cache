// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/cs/cs.proto

package cs

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

type PutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *PutRequest) Reset() {
	*x = PutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cs_cs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutRequest) ProtoMessage() {}

func (x *PutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cs_cs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutRequest.ProtoReflect.Descriptor instead.
func (*PutRequest) Descriptor() ([]byte, []int) {
	return file_proto_cs_cs_proto_rawDescGZIP(), []int{0}
}

func (x *PutRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *PutRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type PutResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *PutResponse) Reset() {
	*x = PutResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_cs_cs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutResponse) ProtoMessage() {}

func (x *PutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cs_cs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutResponse.ProtoReflect.Descriptor instead.
func (*PutResponse) Descriptor() ([]byte, []int) {
	return file_proto_cs_cs_proto_rawDescGZIP(), []int{1}
}

func (x *PutResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_proto_cs_cs_proto protoreflect.FileDescriptor

var file_proto_cs_cs_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x73, 0x2f, 0x63, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x02, 0x63, 0x73, 0x22, 0x34, 0x0a, 0x0a, 0x50, 0x75, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x27, 0x0a,
	0x0b, 0x50, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0x38, 0x0a, 0x0c, 0x43, 0x61, 0x63, 0x68, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x03, 0x50, 0x75, 0x74, 0x12, 0x0e, 0x2e,
	0x63, 0x73, 0x2e, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e,
	0x63, 0x73, 0x2e, 0x50, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x1b, 0x5a, 0x19, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x72, 0x61, 0x66, 0x74, 0x2d, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x63, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_cs_cs_proto_rawDescOnce sync.Once
	file_proto_cs_cs_proto_rawDescData = file_proto_cs_cs_proto_rawDesc
)

func file_proto_cs_cs_proto_rawDescGZIP() []byte {
	file_proto_cs_cs_proto_rawDescOnce.Do(func() {
		file_proto_cs_cs_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_cs_cs_proto_rawDescData)
	})
	return file_proto_cs_cs_proto_rawDescData
}

var file_proto_cs_cs_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_cs_cs_proto_goTypes = []interface{}{
	(*PutRequest)(nil),  // 0: cs.PutRequest
	(*PutResponse)(nil), // 1: cs.PutResponse
}
var file_proto_cs_cs_proto_depIdxs = []int32{
	0, // 0: cs.CacheService.Put:input_type -> cs.PutRequest
	1, // 1: cs.CacheService.Put:output_type -> cs.PutResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_cs_cs_proto_init() }
func file_proto_cs_cs_proto_init() {
	if File_proto_cs_cs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_cs_cs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutRequest); i {
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
		file_proto_cs_cs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutResponse); i {
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
			RawDescriptor: file_proto_cs_cs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_cs_cs_proto_goTypes,
		DependencyIndexes: file_proto_cs_cs_proto_depIdxs,
		MessageInfos:      file_proto_cs_cs_proto_msgTypes,
	}.Build()
	File_proto_cs_cs_proto = out.File
	file_proto_cs_cs_proto_rawDesc = nil
	file_proto_cs_cs_proto_goTypes = nil
	file_proto_cs_cs_proto_depIdxs = nil
}
