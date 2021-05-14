// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.11.4
// source: api/etemplate.proto

package api

import (
	_ "github.com/gogo/protobuf/gogoproto"
	_ "github.com/srikrsna/protoc-gen-gotag/tagger"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Classify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Icon  string `protobuf:"bytes,3,opt,name=icon,proto3" json:"icon,omitempty"`
}

func (x *Classify) Reset() {
	*x = Classify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_etemplate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Classify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Classify) ProtoMessage() {}

func (x *Classify) ProtoReflect() protoreflect.Message {
	mi := &file_api_etemplate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Classify.ProtoReflect.Descriptor instead.
func (*Classify) Descriptor() ([]byte, []int) {
	return file_api_etemplate_proto_rawDescGZIP(), []int{0}
}

func (x *Classify) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Classify) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Classify) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

type UpsertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UpsertRequest) Reset() {
	*x = UpsertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_etemplate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertRequest) ProtoMessage() {}

func (x *UpsertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_etemplate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertRequest.ProtoReflect.Descriptor instead.
func (*UpsertRequest) Descriptor() ([]byte, []int) {
	return file_api_etemplate_proto_rawDescGZIP(), []int{1}
}

func (x *UpsertRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpsertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
}

func (x *UpsertResponse) Reset() {
	*x = UpsertResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_etemplate_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertResponse) ProtoMessage() {}

func (x *UpsertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_etemplate_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertResponse.ProtoReflect.Descriptor instead.
func (*UpsertResponse) Descriptor() ([]byte, []int) {
	return file_api_etemplate_proto_rawDescGZIP(), []int{2}
}

func (x *UpsertResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_api_etemplate_proto protoreflect.FileDescriptor

var file_api_etemplate_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x2d, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x72, 0x69, 0x6b, 0x72, 0x73, 0x6e, 0x61, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x74, 0x61, 0x67, 0x2f, 0x74, 0x61,
	0x67, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x44, 0x0a, 0x08, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x79, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x22, 0x1f, 0x0a, 0x0d, 0x55, 0x70, 0x73, 0x65, 0x72,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x30, 0x0a, 0x0e, 0x55, 0x70, 0x73, 0x65,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0x9a, 0x84, 0x9e, 0x03, 0x09, 0x6a, 0x73, 0x6f,
	0x6e, 0x3a, 0x22, 0x69, 0x64, 0x22, 0x52, 0x02, 0x69, 0x64, 0x32, 0x5e, 0x0a, 0x10, 0x53, 0x75,
	0x62, 0x45, 0x66, 0x66, 0x65, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4a,
	0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55,
	0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x75, 0x62, 0x2d, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x73, 0x42, 0x37, 0x5a, 0x25, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x75, 0x7a, 0x68, 0x6f, 0x6e, 0x67,
	0x7a, 0x68, 0x69, 0x2f, 0x67, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f,
	0x61, 0x70, 0x69, 0xc8, 0xe1, 0x1e, 0x00, 0xd8, 0xe3, 0x1e, 0x00, 0x90, 0xe3, 0x1e, 0x00, 0xd0,
	0xe3, 0x1e, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_etemplate_proto_rawDescOnce sync.Once
	file_api_etemplate_proto_rawDescData = file_api_etemplate_proto_rawDesc
)

func file_api_etemplate_proto_rawDescGZIP() []byte {
	file_api_etemplate_proto_rawDescOnce.Do(func() {
		file_api_etemplate_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_etemplate_proto_rawDescData)
	})
	return file_api_etemplate_proto_rawDescData
}

var file_api_etemplate_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_etemplate_proto_goTypes = []interface{}{
	(*Classify)(nil),       // 0: api.Classify
	(*UpsertRequest)(nil),  // 1: api.UpsertRequest
	(*UpsertResponse)(nil), // 2: api.UpsertResponse
}
var file_api_etemplate_proto_depIdxs = []int32{
	1, // 0: api.SubEffectService.Create:input_type -> api.UpsertRequest
	2, // 1: api.SubEffectService.Create:output_type -> api.UpsertResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_etemplate_proto_init() }
func file_api_etemplate_proto_init() {
	if File_api_etemplate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_etemplate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Classify); i {
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
		file_api_etemplate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertRequest); i {
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
		file_api_etemplate_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertResponse); i {
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
			RawDescriptor: file_api_etemplate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_etemplate_proto_goTypes,
		DependencyIndexes: file_api_etemplate_proto_depIdxs,
		MessageInfos:      file_api_etemplate_proto_msgTypes,
	}.Build()
	File_api_etemplate_proto = out.File
	file_api_etemplate_proto_rawDesc = nil
	file_api_etemplate_proto_goTypes = nil
	file_api_etemplate_proto_depIdxs = nil
}
