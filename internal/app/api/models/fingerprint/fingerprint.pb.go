// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: models/fingerprint/fingerprint.proto

package fingerprint

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

type Fingerprint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Features     []*Fingerprint_Feature `protobuf:"bytes,1,rep,name=features" json:"features,omitempty"`
	Hash         *string                `protobuf:"bytes,2,req,name=hash" json:"hash,omitempty"`
	IsVirtualEnv *bool                  `protobuf:"varint,3,req,name=is_virtual_env,json=isVirtualEnv" json:"is_virtual_env,omitempty"`
	VirtualEnv   []int64                `protobuf:"varint,4,rep,name=virtual_env,json=virtualEnv" json:"virtual_env,omitempty"`
	EnvSources   []int64                `protobuf:"varint,5,rep,name=env_sources,json=envSources" json:"env_sources,omitempty"`
}

func (x *Fingerprint) Reset() {
	*x = Fingerprint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_models_fingerprint_fingerprint_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Fingerprint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Fingerprint) ProtoMessage() {}

func (x *Fingerprint) ProtoReflect() protoreflect.Message {
	mi := &file_models_fingerprint_fingerprint_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Fingerprint.ProtoReflect.Descriptor instead.
func (*Fingerprint) Descriptor() ([]byte, []int) {
	return file_models_fingerprint_fingerprint_proto_rawDescGZIP(), []int{0}
}

func (x *Fingerprint) GetFeatures() []*Fingerprint_Feature {
	if x != nil {
		return x.Features
	}
	return nil
}

func (x *Fingerprint) GetHash() string {
	if x != nil && x.Hash != nil {
		return *x.Hash
	}
	return ""
}

func (x *Fingerprint) GetIsVirtualEnv() bool {
	if x != nil && x.IsVirtualEnv != nil {
		return *x.IsVirtualEnv
	}
	return false
}

func (x *Fingerprint) GetVirtualEnv() []int64 {
	if x != nil {
		return x.VirtualEnv
	}
	return nil
}

func (x *Fingerprint) GetEnvSources() []int64 {
	if x != nil {
		return x.EnvSources
	}
	return nil
}

type Fingerprints struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fps []*Fingerprint `protobuf:"bytes,1,rep,name=fps" json:"fps,omitempty"`
}

func (x *Fingerprints) Reset() {
	*x = Fingerprints{}
	if protoimpl.UnsafeEnabled {
		mi := &file_models_fingerprint_fingerprint_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Fingerprints) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Fingerprints) ProtoMessage() {}

func (x *Fingerprints) ProtoReflect() protoreflect.Message {
	mi := &file_models_fingerprint_fingerprint_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Fingerprints.ProtoReflect.Descriptor instead.
func (*Fingerprints) Descriptor() ([]byte, []int) {
	return file_models_fingerprint_fingerprint_proto_rawDescGZIP(), []int{1}
}

func (x *Fingerprints) GetFps() []*Fingerprint {
	if x != nil {
		return x.Fps
	}
	return nil
}

type Fingerprint_Feature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Feature *int64  `protobuf:"varint,1,req,name=feature" json:"feature,omitempty"`
	Value   *string `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
}

func (x *Fingerprint_Feature) Reset() {
	*x = Fingerprint_Feature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_models_fingerprint_fingerprint_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Fingerprint_Feature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Fingerprint_Feature) ProtoMessage() {}

func (x *Fingerprint_Feature) ProtoReflect() protoreflect.Message {
	mi := &file_models_fingerprint_fingerprint_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Fingerprint_Feature.ProtoReflect.Descriptor instead.
func (*Fingerprint_Feature) Descriptor() ([]byte, []int) {
	return file_models_fingerprint_fingerprint_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Fingerprint_Feature) GetFeature() int64 {
	if x != nil && x.Feature != nil {
		return *x.Feature
	}
	return 0
}

func (x *Fingerprint_Feature) GetValue() string {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return ""
}

var File_models_fingerprint_fingerprint_proto protoreflect.FileDescriptor

var file_models_fingerprint_fingerprint_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x2f, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x22, 0xff, 0x01, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74,
	0x12, 0x39, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x46, 0x69,
	0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x24, 0x0a, 0x0e, 0x69, 0x73, 0x5f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x65, 0x6e,
	0x76, 0x18, 0x03, 0x20, 0x02, 0x28, 0x08, 0x52, 0x0c, 0x69, 0x73, 0x56, 0x69, 0x72, 0x74, 0x75,
	0x61, 0x6c, 0x45, 0x6e, 0x76, 0x12, 0x1f, 0x0a, 0x0b, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c,
	0x5f, 0x65, 0x6e, 0x76, 0x18, 0x04, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0a, 0x76, 0x69, 0x72, 0x74,
	0x75, 0x61, 0x6c, 0x45, 0x6e, 0x76, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6e, 0x76, 0x5f, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x6e, 0x76,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x1a, 0x39, 0x0a, 0x07, 0x46, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20,
	0x02, 0x28, 0x03, 0x52, 0x07, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x37, 0x0a, 0x0c, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e,
	0x74, 0x73, 0x12, 0x27, 0x0a, 0x03, 0x66, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x46, 0x69, 0x6e, 0x67, 0x65,
	0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x52, 0x03, 0x66, 0x70, 0x73, 0x42, 0x16, 0x5a, 0x14, 0x2e,
	0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72,
	0x69, 0x6e, 0x74,
}

var (
	file_models_fingerprint_fingerprint_proto_rawDescOnce sync.Once
	file_models_fingerprint_fingerprint_proto_rawDescData = file_models_fingerprint_fingerprint_proto_rawDesc
)

func file_models_fingerprint_fingerprint_proto_rawDescGZIP() []byte {
	file_models_fingerprint_fingerprint_proto_rawDescOnce.Do(func() {
		file_models_fingerprint_fingerprint_proto_rawDescData = protoimpl.X.CompressGZIP(file_models_fingerprint_fingerprint_proto_rawDescData)
	})
	return file_models_fingerprint_fingerprint_proto_rawDescData
}

var file_models_fingerprint_fingerprint_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_models_fingerprint_fingerprint_proto_goTypes = []interface{}{
	(*Fingerprint)(nil),         // 0: protocol.Fingerprint
	(*Fingerprints)(nil),        // 1: protocol.Fingerprints
	(*Fingerprint_Feature)(nil), // 2: protocol.Fingerprint.Feature
}
var file_models_fingerprint_fingerprint_proto_depIdxs = []int32{
	2, // 0: protocol.Fingerprint.features:type_name -> protocol.Fingerprint.Feature
	0, // 1: protocol.Fingerprints.fps:type_name -> protocol.Fingerprint
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_models_fingerprint_fingerprint_proto_init() }
func file_models_fingerprint_fingerprint_proto_init() {
	if File_models_fingerprint_fingerprint_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_models_fingerprint_fingerprint_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Fingerprint); i {
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
		file_models_fingerprint_fingerprint_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Fingerprints); i {
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
		file_models_fingerprint_fingerprint_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Fingerprint_Feature); i {
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
			RawDescriptor: file_models_fingerprint_fingerprint_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_models_fingerprint_fingerprint_proto_goTypes,
		DependencyIndexes: file_models_fingerprint_fingerprint_proto_depIdxs,
		MessageInfos:      file_models_fingerprint_fingerprint_proto_msgTypes,
	}.Build()
	File_models_fingerprint_fingerprint_proto = out.File
	file_models_fingerprint_fingerprint_proto_rawDesc = nil
	file_models_fingerprint_fingerprint_proto_goTypes = nil
	file_models_fingerprint_fingerprint_proto_depIdxs = nil
}
