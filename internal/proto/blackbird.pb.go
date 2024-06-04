// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: proto/blackbird.proto

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

type GetMetricRequest_Type int32

const (
	GetMetricRequest_CounterMetric GetMetricRequest_Type = 0
	GetMetricRequest_GaugeMetric   GetMetricRequest_Type = 1
)

// Enum value maps for GetMetricRequest_Type.
var (
	GetMetricRequest_Type_name = map[int32]string{
		0: "CounterMetric",
		1: "GaugeMetric",
	}
	GetMetricRequest_Type_value = map[string]int32{
		"CounterMetric": 0,
		"GaugeMetric":   1,
	}
)

func (x GetMetricRequest_Type) Enum() *GetMetricRequest_Type {
	p := new(GetMetricRequest_Type)
	*p = x
	return p
}

func (x GetMetricRequest_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GetMetricRequest_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_blackbird_proto_enumTypes[0].Descriptor()
}

func (GetMetricRequest_Type) Type() protoreflect.EnumType {
	return &file_proto_blackbird_proto_enumTypes[0]
}

func (x GetMetricRequest_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GetMetricRequest_Type.Descriptor instead.
func (GetMetricRequest_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_blackbird_proto_rawDescGZIP(), []int{2, 0}
}

type UpdateMetricRequest_Type int32

const (
	UpdateMetricRequest_CounterMetric UpdateMetricRequest_Type = 0
	UpdateMetricRequest_GaugeMetric   UpdateMetricRequest_Type = 1
)

// Enum value maps for UpdateMetricRequest_Type.
var (
	UpdateMetricRequest_Type_name = map[int32]string{
		0: "CounterMetric",
		1: "GaugeMetric",
	}
	UpdateMetricRequest_Type_value = map[string]int32{
		"CounterMetric": 0,
		"GaugeMetric":   1,
	}
)

func (x UpdateMetricRequest_Type) Enum() *UpdateMetricRequest_Type {
	p := new(UpdateMetricRequest_Type)
	*p = x
	return p
}

func (x UpdateMetricRequest_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UpdateMetricRequest_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_blackbird_proto_enumTypes[1].Descriptor()
}

func (UpdateMetricRequest_Type) Type() protoreflect.EnumType {
	return &file_proto_blackbird_proto_enumTypes[1]
}

func (x UpdateMetricRequest_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UpdateMetricRequest_Type.Descriptor instead.
func (UpdateMetricRequest_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_blackbird_proto_rawDescGZIP(), []int{4, 0}
}

type UpdateMetricsRequest_Type int32

const (
	UpdateMetricsRequest_CounterMetric UpdateMetricsRequest_Type = 0
	UpdateMetricsRequest_GaugeMetric   UpdateMetricsRequest_Type = 1
)

// Enum value maps for UpdateMetricsRequest_Type.
var (
	UpdateMetricsRequest_Type_name = map[int32]string{
		0: "CounterMetric",
		1: "GaugeMetric",
	}
	UpdateMetricsRequest_Type_value = map[string]int32{
		"CounterMetric": 0,
		"GaugeMetric":   1,
	}
)

func (x UpdateMetricsRequest_Type) Enum() *UpdateMetricsRequest_Type {
	p := new(UpdateMetricsRequest_Type)
	*p = x
	return p
}

func (x UpdateMetricsRequest_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UpdateMetricsRequest_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_blackbird_proto_enumTypes[2].Descriptor()
}

func (UpdateMetricsRequest_Type) Type() protoreflect.EnumType {
	return &file_proto_blackbird_proto_enumTypes[2]
}

func (x UpdateMetricsRequest_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UpdateMetricsRequest_Type.Descriptor instead.
func (UpdateMetricsRequest_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_blackbird_proto_rawDescGZIP(), []int{6, 0}
}

type CounterMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Value int64  `protobuf:"varint,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *CounterMetric) Reset() {
	*x = CounterMetric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_blackbird_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CounterMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CounterMetric) ProtoMessage() {}

func (x *CounterMetric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CounterMetric.ProtoReflect.Descriptor instead.
func (*CounterMetric) Descriptor() ([]byte, []int) {
	return file_proto_blackbird_proto_rawDescGZIP(), []int{0}
}

func (x *CounterMetric) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CounterMetric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CounterMetric) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type GaugeMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Value float64 `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GaugeMetric) Reset() {
	*x = GaugeMetric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_blackbird_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GaugeMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GaugeMetric) ProtoMessage() {}

func (x *GaugeMetric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GaugeMetric.ProtoReflect.Descriptor instead.
func (*GaugeMetric) Descriptor() ([]byte, []int) {
	return file_proto_blackbird_proto_rawDescGZIP(), []int{1}
}

func (x *GaugeMetric) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GaugeMetric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GaugeMetric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type GetMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type GetMetricRequest_Type `protobuf:"varint,1,opt,name=type,proto3,enum=main.GetMetricRequest_Type" json:"type,omitempty"`
}

func (x *GetMetricRequest) Reset() {
	*x = GetMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_blackbird_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequest) ProtoMessage() {}

func (x *GetMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[2]
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
	return file_proto_blackbird_proto_rawDescGZIP(), []int{2}
}

func (x *GetMetricRequest) GetType() GetMetricRequest_Type {
	if x != nil {
		return x.Type
	}
	return GetMetricRequest_CounterMetric
}

type GetMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetMetricResponse) Reset() {
	*x = GetMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_blackbird_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponse) ProtoMessage() {}

func (x *GetMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[3]
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
	return file_proto_blackbird_proto_rawDescGZIP(), []int{3}
}

func (x *GetMetricResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
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

	Type UpdateMetricRequest_Type `protobuf:"varint,1,opt,name=type,proto3,enum=main.UpdateMetricRequest_Type" json:"type,omitempty"`
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_blackbird_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[4]
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
	return file_proto_blackbird_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateMetricRequest) GetType() UpdateMetricRequest_Type {
	if x != nil {
		return x.Type
	}
	return UpdateMetricRequest_CounterMetric
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
		mi := &file_proto_blackbird_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricResponse) ProtoMessage() {}

func (x *UpdateMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[5]
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
	return file_proto_blackbird_proto_rawDescGZIP(), []int{5}
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

	Type []UpdateMetricsRequest_Type `protobuf:"varint,1,rep,packed,name=type,proto3,enum=main.UpdateMetricsRequest_Type" json:"type,omitempty"`
}

func (x *UpdateMetricsRequest) Reset() {
	*x = UpdateMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_blackbird_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsRequest) ProtoMessage() {}

func (x *UpdateMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[6]
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
	return file_proto_blackbird_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateMetricsRequest) GetType() []UpdateMetricsRequest_Type {
	if x != nil {
		return x.Type
	}
	return nil
}

type ListMetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Counters []*CounterMetric `protobuf:"bytes,1,rep,name=counters,proto3" json:"counters,omitempty"`
	Gauges   []*GaugeMetric   `protobuf:"bytes,2,rep,name=gauges,proto3" json:"gauges,omitempty"`
	Error    string           `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ListMetricsResponse) Reset() {
	*x = ListMetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_blackbird_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetricsResponse) ProtoMessage() {}

func (x *ListMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_blackbird_proto_msgTypes[7]
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
	return file_proto_blackbird_proto_rawDescGZIP(), []int{7}
}

func (x *ListMetricsResponse) GetCounters() []*CounterMetric {
	if x != nil {
		return x.Counters
	}
	return nil
}

func (x *ListMetricsResponse) GetGauges() []*GaugeMetric {
	if x != nil {
		return x.Gauges
	}
	return nil
}

func (x *ListMetricsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_blackbird_proto protoreflect.FileDescriptor

var file_proto_blackbird_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x6c, 0x61, 0x63, 0x6b, 0x62, 0x69, 0x72,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69, 0x6e, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a, 0x0d, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x47, 0x0a, 0x0b, 0x47, 0x61, 0x75, 0x67, 0x65, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x6f,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1b, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x22, 0x2a, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x10, 0x00, 0x12, 0x0f,
	0x0a, 0x0b, 0x47, 0x61, 0x75, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x10, 0x01, 0x22,
	0x3f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x22, 0x75, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x2a, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x61, 0x75, 0x67, 0x65, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x10, 0x01, 0x22, 0x2c, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x77, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x22, 0x2a, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x10, 0x00, 0x12, 0x0f, 0x0a,
	0x0b, 0x47, 0x61, 0x75, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x10, 0x01, 0x22, 0x87,
	0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x08, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x73, 0x12, 0x29, 0x0a, 0x06, 0x67, 0x61, 0x75, 0x67, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47,
	0x61, 0x75, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x67, 0x61, 0x75, 0x67,
	0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x9c, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x3c, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x12, 0x16, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x45, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x12, 0x19, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0d, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1a, 0x2e, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x43, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x6c, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x6d,
	0x61, 0x69, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x62, 0x61, 0x73, 0x74, 0x74, 0x69, 0x61, 0x6e,
	0x6f, 0x2f, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x62, 0x69, 0x72, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_blackbird_proto_rawDescOnce sync.Once
	file_proto_blackbird_proto_rawDescData = file_proto_blackbird_proto_rawDesc
)

func file_proto_blackbird_proto_rawDescGZIP() []byte {
	file_proto_blackbird_proto_rawDescOnce.Do(func() {
		file_proto_blackbird_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_blackbird_proto_rawDescData)
	})
	return file_proto_blackbird_proto_rawDescData
}

var file_proto_blackbird_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_proto_blackbird_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_blackbird_proto_goTypes = []interface{}{
	(GetMetricRequest_Type)(0),     // 0: main.GetMetricRequest.Type
	(UpdateMetricRequest_Type)(0),  // 1: main.UpdateMetricRequest.Type
	(UpdateMetricsRequest_Type)(0), // 2: main.UpdateMetricsRequest.Type
	(*CounterMetric)(nil),          // 3: main.CounterMetric
	(*GaugeMetric)(nil),            // 4: main.GaugeMetric
	(*GetMetricRequest)(nil),       // 5: main.GetMetricRequest
	(*GetMetricResponse)(nil),      // 6: main.GetMetricResponse
	(*UpdateMetricRequest)(nil),    // 7: main.UpdateMetricRequest
	(*UpdateMetricResponse)(nil),   // 8: main.UpdateMetricResponse
	(*UpdateMetricsRequest)(nil),   // 9: main.UpdateMetricsRequest
	(*ListMetricsResponse)(nil),    // 10: main.ListMetricsResponse
	(*emptypb.Empty)(nil),          // 11: google.protobuf.Empty
}
var file_proto_blackbird_proto_depIdxs = []int32{
	0,  // 0: main.GetMetricRequest.type:type_name -> main.GetMetricRequest.Type
	1,  // 1: main.UpdateMetricRequest.type:type_name -> main.UpdateMetricRequest.Type
	2,  // 2: main.UpdateMetricsRequest.type:type_name -> main.UpdateMetricsRequest.Type
	3,  // 3: main.ListMetricsResponse.counters:type_name -> main.CounterMetric
	4,  // 4: main.ListMetricsResponse.gauges:type_name -> main.GaugeMetric
	5,  // 5: main.Metrics.GetMetric:input_type -> main.GetMetricRequest
	7,  // 6: main.Metrics.UpdateMetric:input_type -> main.UpdateMetricRequest
	9,  // 7: main.Metrics.UpdateMetrics:input_type -> main.UpdateMetricsRequest
	11, // 8: main.Metrics.ListAllMetrics:input_type -> google.protobuf.Empty
	6,  // 9: main.Metrics.GetMetric:output_type -> main.GetMetricResponse
	8,  // 10: main.Metrics.UpdateMetric:output_type -> main.UpdateMetricResponse
	8,  // 11: main.Metrics.UpdateMetrics:output_type -> main.UpdateMetricResponse
	10, // 12: main.Metrics.ListAllMetrics:output_type -> main.ListMetricsResponse
	9,  // [9:13] is the sub-list for method output_type
	5,  // [5:9] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_proto_blackbird_proto_init() }
func file_proto_blackbird_proto_init() {
	if File_proto_blackbird_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_blackbird_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CounterMetric); i {
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
		file_proto_blackbird_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GaugeMetric); i {
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
		file_proto_blackbird_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_blackbird_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_blackbird_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_blackbird_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_blackbird_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_blackbird_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_proto_blackbird_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_blackbird_proto_goTypes,
		DependencyIndexes: file_proto_blackbird_proto_depIdxs,
		EnumInfos:         file_proto_blackbird_proto_enumTypes,
		MessageInfos:      file_proto_blackbird_proto_msgTypes,
	}.Build()
	File_proto_blackbird_proto = out.File
	file_proto_blackbird_proto_rawDesc = nil
	file_proto_blackbird_proto_goTypes = nil
	file_proto_blackbird_proto_depIdxs = nil
}
