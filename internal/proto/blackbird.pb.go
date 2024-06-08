// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: internal/proto/blackbird.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetricType int32

const (
	MetricType_counter MetricType = 0
	MetricType_gauge   MetricType = 1
)

// Enum value maps for MetricType.
var (
	MetricType_name = map[int32]string{
		0: "counter",
		1: "gauge",
	}
	MetricType_value = map[string]int32{
		"counter": 0,
		"gauge":   1,
	}
)

func (x MetricType) Enum() *MetricType {
	p := new(MetricType)
	*p = x
	return p
}

func (x MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_proto_blackbird_proto_enumTypes[0].Descriptor()
}

func (MetricType) Type() protoreflect.EnumType {
	return &file_internal_proto_blackbird_proto_enumTypes[0]
}

func (x MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetricType.Descriptor instead.
func (MetricType) EnumDescriptor() ([]byte, []int) {
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Delta int64      `protobuf:"varint,2,opt,name=delta,proto3" json:"delta,omitempty"`
	Value float64    `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
	Type  MetricType `protobuf:"varint,4,opt,name=type,proto3,enum=main.MetricType" json:"type,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_blackbird_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_blackbird_proto_msgTypes[0]
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
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetType() MetricType {
	if x != nil {
		return x.Type
	}
	return MetricType_counter
}

type GetMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *GetMetricRequest) Reset() {
	*x = GetMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_blackbird_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequest) ProtoMessage() {}

func (x *GetMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_blackbird_proto_msgTypes[1]
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
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{1}
}

func (x *GetMetricRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
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
		mi := &file_internal_proto_blackbird_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponse) ProtoMessage() {}

func (x *GetMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_blackbird_proto_msgTypes[2]
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
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{2}
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

type UpdateMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Value string     `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Type  MetricType `protobuf:"varint,3,opt,name=type,proto3,enum=main.MetricType" json:"type,omitempty"`
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_blackbird_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_blackbird_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateMetricRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateMetricRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *UpdateMetricRequest) GetType() MetricType {
	if x != nil {
		return x.Type
	}
	return MetricType_counter
}

type UpdateMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdateMetricResponse) Reset() {
	*x = UpdateMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_blackbird_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricResponse) ProtoMessage() {}

func (x *UpdateMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_blackbird_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetricResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdateMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *UpdateMetricsRequest) Reset() {
	*x = UpdateMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_blackbird_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsRequest) ProtoMessage() {}

func (x *UpdateMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_blackbird_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricsRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateMetricsRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
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
		mi := &file_internal_proto_blackbird_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetricsResponse) ProtoMessage() {}

func (x *ListMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_blackbird_proto_msgTypes[6]
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
	return file_internal_proto_blackbird_proto_rawDescGZIP(), []int{6}
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

var File_internal_proto_blackbird_proto protoreflect.FileDescriptor

var file_internal_proto_blackbird_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x62, 0x6c, 0x61, 0x63, 0x6b, 0x62, 0x69, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x6d, 0x61, 0x69, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x6a, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65,
	0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22,
	0x38, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x4f, 0x0a, 0x11, 0x47, 0x65, 0x74,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x61, 0x0a, 0x13, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x2c, 0x0a,
	0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3e, 0x0a, 0x14, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x53, 0x0a, 0x13, 0x4c,
	0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x26, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x2a, 0x24, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b,
	0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x67,
	0x61, 0x75, 0x67, 0x65, 0x10, 0x01, 0x32, 0x9c, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x12, 0x3c, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12,
	0x16, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47,
	0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x45, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x12, 0x19, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x43, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x62, 0x61, 0x73, 0x74, 0x74, 0x69, 0x61, 0x6e, 0x6f, 0x2f,
	0x42, 0x6c, 0x61, 0x63, 0x6b, 0x62, 0x69, 0x72, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_blackbird_proto_rawDescOnce sync.Once
	file_internal_proto_blackbird_proto_rawDescData = file_internal_proto_blackbird_proto_rawDesc
)

func file_internal_proto_blackbird_proto_rawDescGZIP() []byte {
	file_internal_proto_blackbird_proto_rawDescOnce.Do(func() {
		file_internal_proto_blackbird_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_blackbird_proto_rawDescData)
	})
	return file_internal_proto_blackbird_proto_rawDescData
}

var file_internal_proto_blackbird_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_proto_blackbird_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_proto_blackbird_proto_goTypes = []interface{}{
	(MetricType)(0),              // 0: main.MetricType
	(*Metric)(nil),               // 1: main.Metric
	(*GetMetricRequest)(nil),     // 2: main.GetMetricRequest
	(*GetMetricResponse)(nil),    // 3: main.GetMetricResponse
	(*UpdateMetricRequest)(nil),  // 4: main.UpdateMetricRequest
	(*UpdateMetricResponse)(nil), // 5: main.UpdateMetricResponse
	(*UpdateMetricsRequest)(nil), // 6: main.UpdateMetricsRequest
	(*ListMetricsResponse)(nil),  // 7: main.ListMetricsResponse
	(*emptypb.Empty)(nil),        // 8: google.protobuf.Empty
}
var file_internal_proto_blackbird_proto_depIdxs = []int32{
	0,  // 0: main.Metric.type:type_name -> main.MetricType
	1,  // 1: main.GetMetricRequest.metric:type_name -> main.Metric
	1,  // 2: main.GetMetricResponse.metric:type_name -> main.Metric
	0,  // 3: main.UpdateMetricRequest.type:type_name -> main.MetricType
	1,  // 4: main.UpdateMetricsRequest.metrics:type_name -> main.Metric
	1,  // 5: main.ListMetricsResponse.metrics:type_name -> main.Metric
	2,  // 6: main.Metrics.GetMetric:input_type -> main.GetMetricRequest
	4,  // 7: main.Metrics.UpdateMetric:input_type -> main.UpdateMetricRequest
	6,  // 8: main.Metrics.UpdateMetrics:input_type -> main.UpdateMetricsRequest
	8,  // 9: main.Metrics.ListAllMetrics:input_type -> google.protobuf.Empty
	3,  // 10: main.Metrics.GetMetric:output_type -> main.GetMetricResponse
	5,  // 11: main.Metrics.UpdateMetric:output_type -> main.UpdateMetricResponse
	5,  // 12: main.Metrics.UpdateMetrics:output_type -> main.UpdateMetricResponse
	7,  // 13: main.Metrics.ListAllMetrics:output_type -> main.ListMetricsResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_internal_proto_blackbird_proto_init() }
func file_internal_proto_blackbird_proto_init() {
	if File_internal_proto_blackbird_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_blackbird_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_proto_blackbird_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_proto_blackbird_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_proto_blackbird_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricRequest); i {
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
		file_internal_proto_blackbird_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricResponse); i {
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
		file_internal_proto_blackbird_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricsRequest); i {
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
		file_internal_proto_blackbird_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_internal_proto_blackbird_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_blackbird_proto_goTypes,
		DependencyIndexes: file_internal_proto_blackbird_proto_depIdxs,
		EnumInfos:         file_internal_proto_blackbird_proto_enumTypes,
		MessageInfos:      file_internal_proto_blackbird_proto_msgTypes,
	}.Build()
	File_internal_proto_blackbird_proto = out.File
	file_internal_proto_blackbird_proto_rawDesc = nil
	file_internal_proto_blackbird_proto_goTypes = nil
	file_internal_proto_blackbird_proto_depIdxs = nil
}
