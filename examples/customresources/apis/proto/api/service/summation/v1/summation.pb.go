// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: api/service/summation/v1/summation.proto

package v1

import (
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

type SumRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Numbers []float64 `protobuf:"fixed64,2,rep,packed,name=numbers,proto3" json:"numbers,omitempty"`
}

func (x *SumRequest) Reset() {
	*x = SumRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_service_summation_v1_summation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SumRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SumRequest) ProtoMessage() {}

func (x *SumRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_service_summation_v1_summation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SumRequest.ProtoReflect.Descriptor instead.
func (*SumRequest) Descriptor() ([]byte, []int) {
	return file_api_service_summation_v1_summation_proto_rawDescGZIP(), []int{0}
}

func (x *SumRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SumRequest) GetNumbers() []float64 {
	if x != nil {
		return x.Numbers
	}
	return nil
}

type SumResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sum float64 `protobuf:"fixed64,1,opt,name=sum,proto3" json:"sum,omitempty"`
}

func (x *SumResponse) Reset() {
	*x = SumResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_service_summation_v1_summation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SumResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SumResponse) ProtoMessage() {}

func (x *SumResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_service_summation_v1_summation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SumResponse.ProtoReflect.Descriptor instead.
func (*SumResponse) Descriptor() ([]byte, []int) {
	return file_api_service_summation_v1_summation_proto_rawDescGZIP(), []int{1}
}

func (x *SumResponse) GetSum() float64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

var File_api_service_summation_v1_summation_proto protoreflect.FileDescriptor

var file_api_service_summation_v1_summation_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73, 0x75,
	0x6d, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x75, 0x6d, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x61, 0x63, 0x6d, 0x65,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x0a, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x01, 0x52, 0x07, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x22,
	0x1f, 0x0a, 0x0b, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x73, 0x75, 0x6d,
	0x32, 0x9c, 0x01, 0x0a, 0x10, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x87, 0x01, 0x0a, 0x03, 0x53, 0x75, 0x6d, 0x12, 0x25, 0x2e,
	0x61, 0x63, 0x6d, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x73, 0x75, 0x6d,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x2b, 0x22, 0x29, 0x2f, 0x61, 0x63, 0x6d, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73, 0x75, 0x6d, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x2f, 0x73, 0x75, 0x6d, 0x42,
	0x2c, 0x5a, 0x2a, 0x67, 0x6f, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_service_summation_v1_summation_proto_rawDescOnce sync.Once
	file_api_service_summation_v1_summation_proto_rawDescData = file_api_service_summation_v1_summation_proto_rawDesc
)

func file_api_service_summation_v1_summation_proto_rawDescGZIP() []byte {
	file_api_service_summation_v1_summation_proto_rawDescOnce.Do(func() {
		file_api_service_summation_v1_summation_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_service_summation_v1_summation_proto_rawDescData)
	})
	return file_api_service_summation_v1_summation_proto_rawDescData
}

var file_api_service_summation_v1_summation_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_service_summation_v1_summation_proto_goTypes = []interface{}{
	(*SumRequest)(nil),  // 0: acme.service.summation.v1.SumRequest
	(*SumResponse)(nil), // 1: acme.service.summation.v1.SumResponse
}
var file_api_service_summation_v1_summation_proto_depIdxs = []int32{
	0, // 0: acme.service.summation.v1.SummationService.Sum:input_type -> acme.service.summation.v1.SumRequest
	1, // 1: acme.service.summation.v1.SummationService.Sum:output_type -> acme.service.summation.v1.SumResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_service_summation_v1_summation_proto_init() }
func file_api_service_summation_v1_summation_proto_init() {
	if File_api_service_summation_v1_summation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_service_summation_v1_summation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SumRequest); i {
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
		file_api_service_summation_v1_summation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SumResponse); i {
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
			RawDescriptor: file_api_service_summation_v1_summation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_service_summation_v1_summation_proto_goTypes,
		DependencyIndexes: file_api_service_summation_v1_summation_proto_depIdxs,
		MessageInfos:      file_api_service_summation_v1_summation_proto_msgTypes,
	}.Build()
	File_api_service_summation_v1_summation_proto = out.File
	file_api_service_summation_v1_summation_proto_rawDesc = nil
	file_api_service_summation_v1_summation_proto_goTypes = nil
	file_api_service_summation_v1_summation_proto_depIdxs = nil
}
