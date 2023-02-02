/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gwwrapper

import (
	"github.com/WesleyWu/gowing/protobuf/gwtypes"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func WrapDouble(v float64) *float64 {
	return &v
}

func WrapFloat(v float32) *float32 {
	return &v
}

func WrapInt64(v int64) *int64 {
	return &v
}

func WrapUInt64(v uint64) *uint64 {
	return &v
}

func WrapInt32(v int32) *int32 {
	return &v
}

func WrapUInt32(v uint32) *uint32 {
	return &v
}

func WrapBool(v bool) *bool {
	return &v
}

func WrapString(v string) *string {
	return &v
}

func AnyDouble(v float64) *anypb.Any {
	valueAny := &wrapperspb.DoubleValue{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyFloat(v float32) *anypb.Any {
	valueAny := &wrapperspb.FloatValue{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyInt64(v int64) *anypb.Any {
	valueAny := &wrapperspb.Int64Value{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyUInt64(v uint64) *anypb.Any {
	valueAny := &wrapperspb.UInt64Value{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyInt32(v int32) *anypb.Any {
	valueAny := &wrapperspb.Int32Value{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyUInt32(v uint32) *anypb.Any {
	valueAny := &wrapperspb.UInt32Value{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyBool(v bool) *anypb.Any {
	valueAny := &wrapperspb.BoolValue{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyString(v string) *anypb.Any {
	valueAny := &wrapperspb.StringValue{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyDoubleSlice(v []float64) *anypb.Any {
	valueAny := &gwtypes.DoubleSlice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyFloatSlice(v []float32) *anypb.Any {
	valueAny := &gwtypes.FloatSlice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyInt64Slice(v []int64) *anypb.Any {
	valueAny := &gwtypes.Int64Slice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyUInt64Slice(v []uint64) *anypb.Any {
	valueAny := &gwtypes.UInt64Slice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyInt32Slice(v []int32) *anypb.Any {
	valueAny := &gwtypes.Int32Slice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyUInt32Slice(v []uint32) *anypb.Any {
	valueAny := &gwtypes.UInt32Slice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyBoolSlice(v []bool) *anypb.Any {
	valueAny := &gwtypes.BoolSlice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyStringSlice(v []string) *anypb.Any {
	valueAny := &gwtypes.StringSlice{Value: v}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyCondition(operator gwtypes.OperatorType, multi gwtypes.MultiType, wildcard gwtypes.WildcardType, value *anypb.Any) *anypb.Any {
	valueCondition := &gwtypes.Condition{
		Operator: operator,
		Multi:    multi,
		Wildcard: wildcard,
		Value:    value,
	}
	result, _ := anypb.New(valueCondition)
	return result
}
