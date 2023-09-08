//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: gwtypes/gwtypes.proto

package gwtypes

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OperatorType int32

const (
	OperatorType_EQ      OperatorType = 0
	OperatorType_NE      OperatorType = 1
	OperatorType_GT      OperatorType = 2
	OperatorType_GTE     OperatorType = 3
	OperatorType_LT      OperatorType = 4
	OperatorType_LTE     OperatorType = 5
	OperatorType_Like    OperatorType = 6
	OperatorType_NotLike OperatorType = 7
	OperatorType_Null    OperatorType = 8
	OperatorType_NotNull OperatorType = 9
)

// Enum value maps for OperatorType.
var (
	OperatorType_name = map[int32]string{
		0: "EQ",
		1: "NE",
		2: "GT",
		3: "GTE",
		4: "LT",
		5: "LTE",
		6: "Like",
		7: "NotLike",
		8: "Null",
		9: "NotNull",
	}
	OperatorType_value = map[string]int32{
		"EQ":      0,
		"NE":      1,
		"GT":      2,
		"GTE":     3,
		"LT":      4,
		"LTE":     5,
		"Like":    6,
		"NotLike": 7,
		"Null":    8,
		"NotNull": 9,
	}
)

func (x OperatorType) Enum() *OperatorType {
	p := new(OperatorType)
	*p = x
	return p
}

func (x OperatorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OperatorType) Descriptor() protoreflect.EnumDescriptor {
	return file_gwtypes_gwtypes_proto_enumTypes[0].Descriptor()
}

func (OperatorType) Type() protoreflect.EnumType {
	return &file_gwtypes_gwtypes_proto_enumTypes[0]
}

func (x OperatorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OperatorType.Descriptor instead.
func (OperatorType) EnumDescriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{0}
}

type MultiType int32

const (
	MultiType_Exact      MultiType = 0
	MultiType_Between    MultiType = 1
	MultiType_NotBetween MultiType = 2
	MultiType_In         MultiType = 3
	MultiType_NotIn      MultiType = 4
)

// Enum value maps for MultiType.
var (
	MultiType_name = map[int32]string{
		0: "Exact",
		1: "Between",
		2: "NotBetween",
		3: "In",
		4: "NotIn",
	}
	MultiType_value = map[string]int32{
		"Exact":      0,
		"Between":    1,
		"NotBetween": 2,
		"In":         3,
		"NotIn":      4,
	}
)

func (x MultiType) Enum() *MultiType {
	p := new(MultiType)
	*p = x
	return p
}

func (x MultiType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MultiType) Descriptor() protoreflect.EnumDescriptor {
	return file_gwtypes_gwtypes_proto_enumTypes[1].Descriptor()
}

func (MultiType) Type() protoreflect.EnumType {
	return &file_gwtypes_gwtypes_proto_enumTypes[1]
}

func (x MultiType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MultiType.Descriptor instead.
func (MultiType) EnumDescriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{1}
}

type WildcardType int32

const (
	WildcardType_None       WildcardType = 0
	WildcardType_Contains   WildcardType = 1
	WildcardType_StartsWith WildcardType = 2
	WildcardType_EndsWith   WildcardType = 3
)

// Enum value maps for WildcardType.
var (
	WildcardType_name = map[int32]string{
		0: "None",
		1: "Contains",
		2: "StartsWith",
		3: "EndsWith",
	}
	WildcardType_value = map[string]int32{
		"None":       0,
		"Contains":   1,
		"StartsWith": 2,
		"EndsWith":   3,
	}
)

func (x WildcardType) Enum() *WildcardType {
	p := new(WildcardType)
	*p = x
	return p
}

func (x WildcardType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WildcardType) Descriptor() protoreflect.EnumDescriptor {
	return file_gwtypes_gwtypes_proto_enumTypes[2].Descriptor()
}

func (WildcardType) Type() protoreflect.EnumType {
	return &file_gwtypes_gwtypes_proto_enumTypes[2]
}

func (x WildcardType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WildcardType.Descriptor instead.
func (WildcardType) EnumDescriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{2}
}

// Wrapper message for array of `double`.
//
// The JSON representation for `DoubleSlice` is JSON array of numbers.
type DoubleSlice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The double value.
	Value []float64 `protobuf:"fixed64,1,rep,packed,name=value,proto3" json:"value,omitempty"`
}

func (x *DoubleSlice) Reset() {
	*x = DoubleSlice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoubleSlice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoubleSlice) ProtoMessage() {}

func (x *DoubleSlice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoubleSlice.ProtoReflect.Descriptor instead.
func (*DoubleSlice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{0}
}

func (x *DoubleSlice) GetValue() []float64 {
	if x != nil {
		return x.Value
	}
	return nil
}

// Wrapper message for array of `float`.
//
// The JSON representation for `FloatSlice` is JSON array of numbers.
type FloatSlice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The float value.
	Value []float32 `protobuf:"fixed32,1,rep,packed,name=value,proto3" json:"value,omitempty"`
}

func (x *FloatSlice) Reset() {
	*x = FloatSlice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FloatSlice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FloatSlice) ProtoMessage() {}

func (x *FloatSlice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FloatSlice.ProtoReflect.Descriptor instead.
func (*FloatSlice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{1}
}

func (x *FloatSlice) GetValue() []float32 {
	if x != nil {
		return x.Value
	}
	return nil
}

// Wrapper message for array of `int64`.
//
// The JSON representation for `Int64Slice` is JSON array of numbers.
type Int64Slice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The int64 value.
	Value []int64 `protobuf:"varint,1,rep,packed,name=value,proto3" json:"value,omitempty"`
}

func (x *Int64Slice) Reset() {
	*x = Int64Slice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Int64Slice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Int64Slice) ProtoMessage() {}

func (x *Int64Slice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Int64Slice.ProtoReflect.Descriptor instead.
func (*Int64Slice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{2}
}

func (x *Int64Slice) GetValue() []int64 {
	if x != nil {
		return x.Value
	}
	return nil
}

// Wrapper message for array of `uint64`.
//
// The JSON representation for `UInt64Slice` is JSON array of numbers.
type UInt64Slice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The uint64 value.
	Value []uint64 `protobuf:"varint,1,rep,packed,name=value,proto3" json:"value,omitempty"`
}

func (x *UInt64Slice) Reset() {
	*x = UInt64Slice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UInt64Slice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UInt64Slice) ProtoMessage() {}

func (x *UInt64Slice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UInt64Slice.ProtoReflect.Descriptor instead.
func (*UInt64Slice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{3}
}

func (x *UInt64Slice) GetValue() []uint64 {
	if x != nil {
		return x.Value
	}
	return nil
}

// Wrapper message for array of `int32`.
//
// The JSON representation for `Int32Slice` is JSON array of numbers.
type Int32Slice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The int32 value.
	Value []int32 `protobuf:"varint,1,rep,packed,name=value,proto3" json:"value,omitempty"`
}

func (x *Int32Slice) Reset() {
	*x = Int32Slice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Int32Slice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Int32Slice) ProtoMessage() {}

func (x *Int32Slice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Int32Slice.ProtoReflect.Descriptor instead.
func (*Int32Slice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{4}
}

func (x *Int32Slice) GetValue() []int32 {
	if x != nil {
		return x.Value
	}
	return nil
}

// Wrapper message for array of `uint32`.
//
// The JSON representation for `UInt32Slice` is JSON array of numbers.
type UInt32Slice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The uint32 value.
	Value []uint32 `protobuf:"varint,1,rep,packed,name=value,proto3" json:"value,omitempty"`
}

func (x *UInt32Slice) Reset() {
	*x = UInt32Slice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UInt32Slice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UInt32Slice) ProtoMessage() {}

func (x *UInt32Slice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UInt32Slice.ProtoReflect.Descriptor instead.
func (*UInt32Slice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{5}
}

func (x *UInt32Slice) GetValue() []uint32 {
	if x != nil {
		return x.Value
	}
	return nil
}

// Wrapper message for array of `bool`.
//
// The JSON representation for `BoolSlice` is JSON array of `true` and `false` strings.
type BoolSlice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The bool value.
	Value []bool `protobuf:"varint,1,rep,packed,name=value,proto3" json:"value,omitempty"`
}

func (x *BoolSlice) Reset() {
	*x = BoolSlice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolSlice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolSlice) ProtoMessage() {}

func (x *BoolSlice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolSlice.ProtoReflect.Descriptor instead.
func (*BoolSlice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{6}
}

func (x *BoolSlice) GetValue() []bool {
	if x != nil {
		return x.Value
	}
	return nil
}

// Wrapper message for array of `string`.
//
// The JSON representation for `StringSlice` is JSON array of strings.
type StringSlice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The string value.
	Value []string `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *StringSlice) Reset() {
	*x = StringSlice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringSlice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringSlice) ProtoMessage() {}

func (x *StringSlice) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringSlice.ProtoReflect.Descriptor instead.
func (*StringSlice) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{7}
}

func (x *StringSlice) GetValue() []string {
	if x != nil {
		return x.Value
	}
	return nil
}

type Condition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operator OperatorType `protobuf:"varint,1,opt,name=operator,proto3,enum=gwtypes.OperatorType" json:"operator,omitempty"`
	Multi    MultiType    `protobuf:"varint,2,opt,name=multi,proto3,enum=gwtypes.MultiType" json:"multi,omitempty"`
	Wildcard WildcardType `protobuf:"varint,3,opt,name=wildcard,proto3,enum=gwtypes.WildcardType" json:"wildcard,omitempty"`
	Value    *anypb.Any   `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Condition) Reset() {
	*x = Condition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gwtypes_gwtypes_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Condition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Condition) ProtoMessage() {}

func (x *Condition) ProtoReflect() protoreflect.Message {
	mi := &file_gwtypes_gwtypes_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Condition.ProtoReflect.Descriptor instead.
func (*Condition) Descriptor() ([]byte, []int) {
	return file_gwtypes_gwtypes_proto_rawDescGZIP(), []int{8}
}

func (x *Condition) GetOperator() OperatorType {
	if x != nil {
		return x.Operator
	}
	return OperatorType_EQ
}

func (x *Condition) GetMulti() MultiType {
	if x != nil {
		return x.Multi
	}
	return MultiType_Exact
}

func (x *Condition) GetWildcard() WildcardType {
	if x != nil {
		return x.Wildcard
	}
	return WildcardType_None
}

func (x *Condition) GetValue() *anypb.Any {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_gwtypes_gwtypes_proto protoreflect.FileDescriptor

var file_gwtypes_gwtypes_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x67, 0x77, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x67, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x23, 0x0a, 0x0b, 0x44,
	0x6f, 0x75, 0x62, 0x6c, 0x65, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x22, 0x0a, 0x0a, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x02, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x53, 0x6c, 0x69,
	0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x55, 0x49, 0x6e, 0x74,
	0x36, 0x34, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x22, 0x0a,
	0x0a, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x23, 0x0a, 0x0b, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x53, 0x6c, 0x69, 0x63, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x21, 0x0a, 0x09, 0x42, 0x6f, 0x6f, 0x6c, 0x53, 0x6c,
	0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xc7,
	0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x08,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15,
	0x2e, 0x67, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12,
	0x28, 0x0a, 0x05, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12,
	0x2e, 0x67, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x05, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x12, 0x31, 0x0a, 0x08, 0x77, 0x69, 0x6c,
	0x64, 0x63, 0x61, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x67, 0x77,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x57, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x08, 0x77, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64, 0x12, 0x2a, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x6e, 0x0a, 0x0c, 0x4f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x45, 0x51, 0x10, 0x00,
	0x12, 0x06, 0x0a, 0x02, 0x4e, 0x45, 0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x47, 0x54, 0x10, 0x02,
	0x12, 0x07, 0x0a, 0x03, 0x47, 0x54, 0x45, 0x10, 0x03, 0x12, 0x06, 0x0a, 0x02, 0x4c, 0x54, 0x10,
	0x04, 0x12, 0x07, 0x0a, 0x03, 0x4c, 0x54, 0x45, 0x10, 0x05, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x69,
	0x6b, 0x65, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x6f, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x10,
	0x07, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x75, 0x6c, 0x6c, 0x10, 0x08, 0x12, 0x0b, 0x0a, 0x07, 0x4e,
	0x6f, 0x74, 0x4e, 0x75, 0x6c, 0x6c, 0x10, 0x09, 0x2a, 0x46, 0x0a, 0x09, 0x4d, 0x75, 0x6c, 0x74,
	0x69, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x78, 0x61, 0x63, 0x74, 0x10, 0x00,
	0x12, 0x0b, 0x0a, 0x07, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x10, 0x01, 0x12, 0x0e, 0x0a,
	0x0a, 0x4e, 0x6f, 0x74, 0x42, 0x65, 0x74, 0x77, 0x65, 0x65, 0x6e, 0x10, 0x02, 0x12, 0x06, 0x0a,
	0x02, 0x49, 0x6e, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x4e, 0x6f, 0x74, 0x49, 0x6e, 0x10, 0x04,
	0x2a, 0x44, 0x0a, 0x0c, 0x57, 0x69, 0x6c, 0x64, 0x63, 0x61, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x73, 0x57, 0x69, 0x74, 0x68, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x6e, 0x64, 0x73,
	0x57, 0x69, 0x74, 0x68, 0x10, 0x03, 0x42, 0x84, 0x01, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x67,
	0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0x42, 0x0c, 0x47, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x73, 0x6c, 0x65, 0x79, 0x77, 0x75, 0x2f, 0x67, 0x6f, 0x77, 0x69,
	0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x77, 0x74, 0x79,
	0x70, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x47, 0x58, 0x58, 0xaa, 0x02, 0x07, 0x47, 0x77, 0x74, 0x79,
	0x70, 0x65, 0x73, 0xca, 0x02, 0x07, 0x47, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0xe2, 0x02, 0x13,
	0x47, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x07, 0x47, 0x77, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gwtypes_gwtypes_proto_rawDescOnce sync.Once
	file_gwtypes_gwtypes_proto_rawDescData = file_gwtypes_gwtypes_proto_rawDesc
)

func file_gwtypes_gwtypes_proto_rawDescGZIP() []byte {
	file_gwtypes_gwtypes_proto_rawDescOnce.Do(func() {
		file_gwtypes_gwtypes_proto_rawDescData = protoimpl.X.CompressGZIP(file_gwtypes_gwtypes_proto_rawDescData)
	})
	return file_gwtypes_gwtypes_proto_rawDescData
}

var file_gwtypes_gwtypes_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_gwtypes_gwtypes_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_gwtypes_gwtypes_proto_goTypes = []interface{}{
	(OperatorType)(0),   // 0: gwtypes.OperatorType
	(MultiType)(0),      // 1: gwtypes.MultiType
	(WildcardType)(0),   // 2: gwtypes.WildcardType
	(*DoubleSlice)(nil), // 3: gwtypes.DoubleSlice
	(*FloatSlice)(nil),  // 4: gwtypes.FloatSlice
	(*Int64Slice)(nil),  // 5: gwtypes.Int64Slice
	(*UInt64Slice)(nil), // 6: gwtypes.UInt64Slice
	(*Int32Slice)(nil),  // 7: gwtypes.Int32Slice
	(*UInt32Slice)(nil), // 8: gwtypes.UInt32Slice
	(*BoolSlice)(nil),   // 9: gwtypes.BoolSlice
	(*StringSlice)(nil), // 10: gwtypes.StringSlice
	(*Condition)(nil),   // 11: gwtypes.Condition
	(*anypb.Any)(nil),   // 12: google.protobuf.Any
}
var file_gwtypes_gwtypes_proto_depIdxs = []int32{
	0,  // 0: gwtypes.Condition.operator:type_name -> gwtypes.OperatorType
	1,  // 1: gwtypes.Condition.multi:type_name -> gwtypes.MultiType
	2,  // 2: gwtypes.Condition.wildcard:type_name -> gwtypes.WildcardType
	12, // 3: gwtypes.Condition.value:type_name -> google.protobuf.Any
	4,  // [4:4] is the sub-list for method output_type
	4,  // [4:4] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_gwtypes_gwtypes_proto_init() }
func file_gwtypes_gwtypes_proto_init() {
	if File_gwtypes_gwtypes_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gwtypes_gwtypes_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoubleSlice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FloatSlice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Int64Slice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UInt64Slice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Int32Slice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UInt32Slice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolSlice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringSlice); i {
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
		file_gwtypes_gwtypes_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Condition); i {
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
			RawDescriptor: file_gwtypes_gwtypes_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gwtypes_gwtypes_proto_goTypes,
		DependencyIndexes: file_gwtypes_gwtypes_proto_depIdxs,
		EnumInfos:         file_gwtypes_gwtypes_proto_enumTypes,
		MessageInfos:      file_gwtypes_gwtypes_proto_msgTypes,
	}.Build()
	File_gwtypes_gwtypes_proto = out.File
	file_gwtypes_gwtypes_proto_rawDesc = nil
	file_gwtypes_gwtypes_proto_goTypes = nil
	file_gwtypes_gwtypes_proto_depIdxs = nil
}
