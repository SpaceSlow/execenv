// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: proto/execenv.proto

package proto

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

type MType int32

const (
	MType_UNSPECIFIED MType = 0
	MType_COUNTER     MType = 1
	MType_GAUGE       MType = 2
)

// Enum value maps for MType.
var (
	MType_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "COUNTER",
		2: "GAUGE",
	}
	MType_value = map[string]int32{
		"UNSPECIFIED": 0,
		"COUNTER":     1,
		"GAUGE":       2,
	}
)

func (x MType) Enum() *MType {
	p := new(MType)
	*p = x
	return p
}

func (x MType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_execenv_proto_enumTypes[0].Descriptor()
}

func (MType) Type() protoreflect.EnumType {
	return &file_proto_execenv_proto_enumTypes[0]
}

func (x MType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MType.Descriptor instead.
func (MType) EnumDescriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MType MType   `protobuf:"varint,2,opt,name=mType,proto3,enum=execenv.MType" json:"mType,omitempty"`
	Value float64 `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
	Delta int64   `protobuf:"varint,4,opt,name=delta,proto3" json:"delta,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetMType() MType {
	if x != nil {
		return x.MType
	}
	return MType_UNSPECIFIED
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

type AddMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *AddMetricRequest) Reset() {
	*x = AddMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMetricRequest) ProtoMessage() {}

func (x *AddMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMetricRequest.ProtoReflect.Descriptor instead.
func (*AddMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{1}
}

func (x *AddMetricRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type AddMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *AddMetricResponse) Reset() {
	*x = AddMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMetricResponse) ProtoMessage() {}

func (x *AddMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMetricResponse.ProtoReflect.Descriptor instead.
func (*AddMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{2}
}

func (x *AddMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type BatchAddMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *BatchAddMetricRequest) Reset() {
	*x = BatchAddMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchAddMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchAddMetricRequest) ProtoMessage() {}

func (x *BatchAddMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchAddMetricRequest.ProtoReflect.Descriptor instead.
func (*BatchAddMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{3}
}

func (x *BatchAddMetricRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type BatchAddMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *BatchAddMetricResponse) Reset() {
	*x = BatchAddMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchAddMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchAddMetricResponse) ProtoMessage() {}

func (x *BatchAddMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchAddMetricResponse.ProtoReflect.Descriptor instead.
func (*BatchAddMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{4}
}

func (x *BatchAddMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	MType MType  `protobuf:"varint,2,opt,name=mType,proto3,enum=execenv.MType" json:"mType,omitempty"`
}

func (x *GetMetricRequest) Reset() {
	*x = GetMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequest) ProtoMessage() {}

func (x *GetMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricRequest.ProtoReflect.Descriptor instead.
func (*GetMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{5}
}

func (x *GetMetricRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetMetricRequest) GetMType() MType {
	if x != nil {
		return x.MType
	}
	return MType_UNSPECIFIED
}

type GetMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Error  string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetMetricResponse) Reset() {
	*x = GetMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponse) ProtoMessage() {}

func (x *GetMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricResponse.ProtoReflect.Descriptor instead.
func (*GetMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{6}
}

func (x *GetMetricResponse) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *GetMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type ListMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListMetricsRequest) Reset() {
	*x = ListMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetricsRequest) ProtoMessage() {}

func (x *ListMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMetricsRequest.ProtoReflect.Descriptor instead.
func (*ListMetricsRequest) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{7}
}

type ListMetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
	Error   string    `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ListMetricsResponse) Reset() {
	*x = ListMetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_execenv_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetricsResponse) ProtoMessage() {}

func (x *ListMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execenv_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMetricsResponse.ProtoReflect.Descriptor instead.
func (*ListMetricsResponse) Descriptor() ([]byte, []int) {
	return file_proto_execenv_proto_rawDescGZIP(), []int{8}
}

func (x *ListMetricsResponse) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

func (x *ListMetricsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_execenv_proto protoreflect.FileDescriptor

var file_proto_execenv_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x22, 0x6a,
	0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x05, 0x6d, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e,
	0x76, 0x2e, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x22, 0x3b, 0x0a, 0x10, 0x41, 0x64,
	0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x29, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x22, 0x42, 0x0a, 0x15, 0x42, 0x61, 0x74, 0x63, 0x68, 0x41, 0x64, 0x64, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x07, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65,
	0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x2e, 0x0a, 0x16, 0x42, 0x61, 0x74, 0x63, 0x68, 0x41,
	0x64, 0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x48, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x05, 0x6d, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x65, 0x78, 0x65, 0x63,
	0x65, 0x6e, 0x76, 0x2e, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x6d, 0x54, 0x79, 0x70, 0x65,
	0x22, 0x52, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x14, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x56, 0x0a, 0x13, 0x4c, 0x69,
	0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x29, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x2a, 0x30, 0x0a, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x4f, 0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x41, 0x55,
	0x47, 0x45, 0x10, 0x02, 0x32, 0xb4, 0x02, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x19, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x41, 0x64,
	0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1e, 0x2e, 0x65,
	0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x41, 0x64, 0x64, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x65,
	0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x41, 0x64, 0x64, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x19, 0x2e, 0x65, 0x78, 0x65,
	0x63, 0x65, 0x6e, 0x76, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e,
	0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x48, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x12, 0x1b, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e,
	0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x70, 0x61, 0x63, 0x65, 0x53,
	0x6c, 0x6f, 0x77, 0x2f, 0x65, 0x78, 0x65, 0x63, 0x65, 0x6e, 0x76, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_execenv_proto_rawDescOnce sync.Once
	file_proto_execenv_proto_rawDescData = file_proto_execenv_proto_rawDesc
)

func file_proto_execenv_proto_rawDescGZIP() []byte {
	file_proto_execenv_proto_rawDescOnce.Do(func() {
		file_proto_execenv_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_execenv_proto_rawDescData)
	})
	return file_proto_execenv_proto_rawDescData
}

var file_proto_execenv_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_execenv_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_execenv_proto_goTypes = []any{
	(MType)(0),                     // 0: execenv.MType
	(*Metric)(nil),                 // 1: execenv.Metric
	(*AddMetricRequest)(nil),       // 2: execenv.AddMetricRequest
	(*AddMetricResponse)(nil),      // 3: execenv.AddMetricResponse
	(*BatchAddMetricRequest)(nil),  // 4: execenv.BatchAddMetricRequest
	(*BatchAddMetricResponse)(nil), // 5: execenv.BatchAddMetricResponse
	(*GetMetricRequest)(nil),       // 6: execenv.GetMetricRequest
	(*GetMetricResponse)(nil),      // 7: execenv.GetMetricResponse
	(*ListMetricsRequest)(nil),     // 8: execenv.ListMetricsRequest
	(*ListMetricsResponse)(nil),    // 9: execenv.ListMetricsResponse
}
var file_proto_execenv_proto_depIdxs = []int32{
	0,  // 0: execenv.Metric.mType:type_name -> execenv.MType
	1,  // 1: execenv.AddMetricRequest.metric:type_name -> execenv.Metric
	1,  // 2: execenv.BatchAddMetricRequest.metrics:type_name -> execenv.Metric
	0,  // 3: execenv.GetMetricRequest.mType:type_name -> execenv.MType
	1,  // 4: execenv.GetMetricResponse.metric:type_name -> execenv.Metric
	1,  // 5: execenv.ListMetricsResponse.metrics:type_name -> execenv.Metric
	2,  // 6: execenv.MetricService.AddMetric:input_type -> execenv.AddMetricRequest
	4,  // 7: execenv.MetricService.BatchAddMetric:input_type -> execenv.BatchAddMetricRequest
	6,  // 8: execenv.MetricService.GetMetric:input_type -> execenv.GetMetricRequest
	8,  // 9: execenv.MetricService.ListMetrics:input_type -> execenv.ListMetricsRequest
	3,  // 10: execenv.MetricService.AddMetric:output_type -> execenv.AddMetricResponse
	5,  // 11: execenv.MetricService.BatchAddMetric:output_type -> execenv.BatchAddMetricResponse
	7,  // 12: execenv.MetricService.GetMetric:output_type -> execenv.GetMetricResponse
	9,  // 13: execenv.MetricService.ListMetrics:output_type -> execenv.ListMetricsResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_proto_execenv_proto_init() }
func file_proto_execenv_proto_init() {
	if File_proto_execenv_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_execenv_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Metric); i {
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
		file_proto_execenv_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AddMetricRequest); i {
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
		file_proto_execenv_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*AddMetricResponse); i {
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
		file_proto_execenv_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*BatchAddMetricRequest); i {
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
		file_proto_execenv_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*BatchAddMetricResponse); i {
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
		file_proto_execenv_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricRequest); i {
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
		file_proto_execenv_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricResponse); i {
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
		file_proto_execenv_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*ListMetricsRequest); i {
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
		file_proto_execenv_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*ListMetricsResponse); i {
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
			RawDescriptor: file_proto_execenv_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_execenv_proto_goTypes,
		DependencyIndexes: file_proto_execenv_proto_depIdxs,
		EnumInfos:         file_proto_execenv_proto_enumTypes,
		MessageInfos:      file_proto_execenv_proto_msgTypes,
	}.Build()
	File_proto_execenv_proto = out.File
	file_proto_execenv_proto_rawDesc = nil
	file_proto_execenv_proto_goTypes = nil
	file_proto_execenv_proto_depIdxs = nil
}
