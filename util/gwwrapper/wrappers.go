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
	"context"
	"github.com/WesleyWu/gowing/protobuf/gwtypes"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
	"strings"
	"time"
)

const (
	ConditionQueryPrefix = "condition{"
	ConditionQuerySuffix = "}"
	TagNameOperator      = "operator"
	TagNameMulti         = "multi"
	TagNameWildcard      = "wildcard"

	FieldNameOperator = "operator"
	FieldNameMulti    = "multi"
	FieldNameWildcard = "wildcard"
	FieldNameValue    = "value"
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

func WrapTimestamp(v time.Time) *timestamppb.Timestamp {
	return timestamppb.New(v)
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

func AnyInt(v int) *anypb.Any {
	valueAny := &wrapperspb.Int32Value{Value: int32(v)}
	result, _ := anypb.New(valueAny)
	return result
}

func AnyUInt(v uint) *anypb.Any {
	valueAny := &wrapperspb.UInt32Value{Value: uint32(v)}
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

func AnyIntSlice(v []int) *anypb.Any {
	value := make([]int32, len(v))
	for i, v := range v {
		value[i] = int32(v)
	}
	return AnyInt32Slice(value)
}

func AnyUIntSlice(v []uint) *anypb.Any {
	value := make([]uint32, len(v))
	for i, v := range v {
		value[i] = uint32(v)
	}
	return AnyUInt32Slice(value)
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

func AnyInterface(ctx context.Context, any interface{}, tag reflect.StructTag) *anypb.Any {
	switch value := any.(type) {
	case string:
		return AnyString(value)
	case int:
		return AnyInt(value)
	case int8, int16, uint8, uint16:
		g.Log().Warning(ctx, "Cannot convert int8,int16,uint8,uint16 types to *anypb.Any")
		return nil
	case int32:
		return AnyInt32(value)
	case int64:
		return AnyInt64(value)
	case uint:
		return AnyUInt(value)
	case uint32:
		return AnyUInt32(value)
	case uint64:
		return AnyUInt64(value)
	case bool:
		return AnyBool(value)
	case float32:
		return AnyFloat(value)
	case float64:
		return AnyDouble(value)
	default:
		t := reflect.TypeOf(any)
		switch t.Kind() {
		case reflect.Slice, reflect.Array:
			return AnySlice(ctx, any, tag)
		case reflect.Ptr:
			return AnyInterface(ctx, reflect.ValueOf(any).Elem(), tag)
		default:
			g.Log().Warningf(ctx, "Cannot convert type %s to *anypb.Any", t.String())
			return nil
		}
		g.Log().Warningf(ctx, "Cannot convert %s types to *anypb.Any", reflect.TypeOf(any).String())
		return nil
	}
}

func AnySlice(ctx context.Context, any interface{}, tag reflect.StructTag) *anypb.Any {
	switch value := any.(type) {
	case []string:
		return AnyStringSlice(value)
	case []int:
		return AnyIntSlice(value)
	case []int8, int16, []uint8, []uint16:
		g.Log().Warning(ctx, "Cannot convert int8,int16,uint8,uint16 types to *anypb.Any")
		return nil
	case []int32:
		return AnyInt32Slice(value)
	case []int64:
		return AnyInt64Slice(value)
	case []uint:
		return AnyUIntSlice(value)
	case []uint32:
		return AnyUInt32Slice(value)
	case []uint64:
		return AnyUInt64Slice(value)
	case []bool:
		return AnyBoolSlice(value)
	case []float32:
		return AnyFloatSlice(value)
	case []float64:
		return AnyDoubleSlice(value)
	default:
		g.Log().Warningf(ctx, "Cannot convert %s types to *anypb.Any", reflect.TypeOf(any).String())
		return nil
	}
}

func AnySliceCondition(ctx context.Context, any interface{}, tag reflect.StructTag) *anypb.Any {
	operator := gwtypes.MustParseOperatorType(tag.Get(TagNameOperator), gwtypes.OperatorType_EQ)
	multi := gwtypes.MustParseMultiType(tag.Get(TagNameMulti), gwtypes.MultiType_In)
	wildcard := gwtypes.MustParseWildcardType(tag.Get(TagNameWildcard), gwtypes.WildcardType_None)
	switch value := any.(type) {
	case []interface{}:
		array := make([]string, len(value))
		for k, v := range value {
			array[k] = gconv.String(v)
		}
		return AnyCondition(operator, multi, wildcard, AnyStringSlice(array))
	case []string:
		return AnyCondition(operator, multi, wildcard, AnyStringSlice(value))
	case []int:
		return AnyCondition(operator, multi, wildcard, AnyIntSlice(value))
	case []int8, []int16, []uint8, []uint16:
		g.Log().Warning(ctx, "Cannot convert []int8,[]int16,[]uint8,[]uint16 types to *anypb.Any")
		return nil
	case []int32:
		return AnyCondition(operator, multi, wildcard, AnyInt32Slice(value))
	case []int64:
		return AnyCondition(operator, multi, wildcard, AnyInt64Slice(value))
	case []uint:
		return AnyCondition(operator, multi, wildcard, AnyUIntSlice(value))
	case []uint32:
		return AnyCondition(operator, multi, wildcard, AnyUInt32Slice(value))
	case []uint64:
		return AnyCondition(operator, multi, wildcard, AnyUInt64Slice(value))
	case []bool:
		return AnyCondition(operator, multi, wildcard, AnyBoolSlice(value))
	case []float32:
		return AnyCondition(operator, multi, wildcard, AnyFloatSlice(value))
	case []float64:
		return AnyCondition(operator, multi, wildcard, AnyDoubleSlice(value))
	default:
		g.Log().Warningf(ctx, "Cannot convert %s types to *anypb.Any", reflect.TypeOf(any).String())
		return nil
	}
}

func AnyInterfaceCondition(ctx context.Context, any interface{}, tag reflect.StructTag) *anypb.Any {
	operator := gwtypes.MustParseOperatorType(tag.Get(TagNameOperator), gwtypes.OperatorType_EQ)
	multi := gwtypes.MustParseMultiType(tag.Get(TagNameMulti), gwtypes.MultiType_Exact)
	wildcard := gwtypes.MustParseWildcardType(tag.Get(TagNameWildcard), gwtypes.WildcardType_None)
	switch value := any.(type) {
	case nil:
		return nil
	case int:
		return AnyCondition(operator, multi, wildcard, AnyInt(value))
	case int32:
		return AnyCondition(operator, multi, wildcard, AnyInt32(value))
	case int64:
		return AnyCondition(operator, multi, wildcard, AnyInt64(value))
	case uint:
		return AnyCondition(operator, multi, wildcard, AnyUInt(value))
	case uint32:
		return AnyCondition(operator, multi, wildcard, AnyUInt32(value))
	case uint64:
		return AnyCondition(operator, multi, wildcard, AnyUInt64(value))
	case bool:
		return AnyCondition(operator, multi, wildcard, AnyBool(value))
	case float32:
		return AnyCondition(operator, multi, wildcard, AnyFloat(value))
	case float64:
		return AnyCondition(operator, multi, wildcard, AnyDouble(value))
	case string:
		if g.IsEmpty(value) {
			return nil
		}
		if strings.HasPrefix(value, ConditionQueryPrefix) && strings.HasSuffix(value, ConditionQuerySuffix) {
			conditionJson, err := gjson.DecodeToJson(value[9:])
			if err != nil {
				g.Log().Warningf(ctx, "Cannot decode %s to gwtypes.Condition", value)
				return nil
			}
			operatorValue := gwtypes.MustParseOperatorType(conditionJson.Get(FieldNameOperator).String(), gwtypes.OperatorType_EQ)
			multiValue := gwtypes.MustParseMultiType(conditionJson.Get(FieldNameMulti).String(), gwtypes.MultiType_Exact)
			wildcardValue := gwtypes.MustParseWildcardType(conditionJson.Get(FieldNameWildcard).String(), gwtypes.WildcardType_None)
			return AnyCondition(operatorValue, multiValue, wildcardValue, AnyInterface(ctx, conditionJson.Get(FieldNameValue).Interface(), tag))
		}
		switch wildcard {
		case gwtypes.WildcardType_Contains:
			return AnyCondition(operator, multi, wildcard, AnyString("%"+value+"%"))
		case gwtypes.WildcardType_StartsWith:
			return AnyCondition(operator, multi, wildcard, AnyString(value+"%"))
		case gwtypes.WildcardType_EndsWith:
			return AnyCondition(operator, multi, wildcard, AnyString("%"+value))
		default:
			return AnyCondition(operator, multi, wildcard, AnyString(value))
		}
	default:
		t := reflect.TypeOf(any)
		switch t.Kind() {
		case reflect.Slice, reflect.Array:
			return AnySliceCondition(ctx, any, tag)
		case reflect.Ptr:
			return AnyInterfaceCondition(ctx, reflect.ValueOf(any).Elem(), tag)
		default:
			g.Log().Warningf(ctx, "Cannot convert type %s to *anypb.Any", t.String())
			return nil
		}
	}
}
