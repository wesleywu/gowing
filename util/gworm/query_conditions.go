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

package gworm

import (
	"context"
	"reflect"
	"strings"

	"github.com/WesleyWu/gowing/errors/gwerror"
	"github.com/WesleyWu/gowing/protobuf/gwtypes"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// todo gworm package is deprecated
const (
	ConditionQueryPrefix = "condition{"
	ConditionQuerySuffix = "}"
	TagNameMulti         = "multi"
	TagNameWildcard      = "wildcard"
)

type FilterRequest struct {
	PropertyFilters []*PropertyFilter
}

type PropertyFilter struct {
	Property string               `json:"property"`
	Value    interface{}          `json:"value"`
	Operator gwtypes.OperatorType `json:"operator"`
	Multi    gwtypes.MultiType    `json:"multi"`
	Wildcard gwtypes.WildcardType `json:"wildcard"`
}

func (fr *FilterRequest) AddPropertyFilter(f *PropertyFilter) *FilterRequest {
	fr.PropertyFilters = append(fr.PropertyFilters, f)
	return fr
}

func ExtractFilters(ctx context.Context, req interface{}, columnMap map[string]string, mType ModelType) (fr FilterRequest, err error) {
	var f *PropertyFilter
	p := reflect.TypeOf(req)
	if p.Kind() != reflect.Ptr { // 要求传入值必须是个指针
		err = gwerror.NewBadRequestErrorf(req, "服务函数的输入参数必须是结构体指针")
		return
	}
	t := p.Elem()
	//g.Log().Debugf(ctx, "kind of input parameter is %s", t.Name())

	queryValue := reflect.ValueOf(req).Elem()

	// 循环结构体的字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		if fieldName == "Meta" {
			continue
		}
		// Only do converting to public attributes.
		if !field.IsExported() {
			continue
		}
		fieldType := field.Type
		//g.Log().Debugf(ctx, "kind of field \"%s\" is %s", fieldName, field.Type.Kind().String())
		switch fieldType.Kind() {
		case reflect.Ptr:
			columnName, exists := columnMap[fieldName]
			if !exists {
				continue
			}
			//g.Log().Debugf(ctx, "kind of element of field %s is %s", fieldName, fieldElemType.Kind().String())
			anyValue := queryValue.Field(i).Interface().(*anypb.Any)
			if anyValue != nil {
				f, err = unwrapAnyFilter(columnName, field.Tag, anyValue, mType)
				if err != nil {
					return
				}
				if f != nil {
					fr.AddPropertyFilter(f)
				}
			}

		case reflect.Struct:
			structValue := queryValue.Field(i)
			g.Log().Debugf(ctx, "value of field %s is %x", fieldName, structValue)
			for si := 0; si < fieldType.NumField(); si++ {
				innerField := fieldType.Field(si)
				if innerField.Type.Kind() != reflect.Interface { // 仅处理类型为 interface{} 的字段
					continue
				}
				columnName, exists := columnMap[innerField.Name] // 仅处理在表字段定义中有的字段
				if !exists {
					continue
				}
				fieldValue := structValue.Field(si).Interface()
				if fieldValue == nil { // 不出来值为nil的字段
					continue
				}
				g.Log().Debugf(ctx, "inner field %s kind:%si, column:%s, value:%s", innerField.Name, innerField.Type.Kind().String(), columnName, fieldValue)
				f, err = parsePropertyFilter(ctx, req, columnName, innerField.Tag, fieldValue, mType)
				if err != nil {
					return
				}
				fr.AddPropertyFilter(f)
			}
		case reflect.Interface:
			columnName, exists := columnMap[fieldName]
			if !exists {
				continue
			}
			fieldValue := queryValue.Field(i).Interface()
			if fieldValue == nil {
				continue
			}
			f, err = parsePropertyFilter(ctx, req, columnName, field.Tag, fieldValue, mType)
			if err != nil {
				return
			}
			fr.AddPropertyFilter(f)
		}
	}
	return
}

func parsePropertyFilter(ctx context.Context, req interface{}, columnName string, tag reflect.StructTag, value interface{}, mType ModelType) (*PropertyFilter, error) {
	if value == nil { // todo processing: is null/is not null
		return nil, nil
	}
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem() == reflect.TypeOf(PropertyFilter{}) {
			pf := value.(*PropertyFilter)
			pf.Property = columnName
			return pf, nil
		}
		return &PropertyFilter{
			Property: columnName,
			Value:    value,
			Operator: gwtypes.OperatorType_EQ,
			Multi:    gwtypes.MultiType_Exact,
			Wildcard: gwtypes.WildcardType_None,
		}, nil
	case reflect.Slice, reflect.Array:
		valueSlice := gconv.SliceAny(value)
		multiTag, ok := tag.Lookup(TagNameMulti)
		if ok {
			multi, err := gwtypes.ParseMultiType(multiTag)
			if err != nil {
				return nil, err
			}
			switch len(valueSlice) {
			case 0:
				return nil, nil
			case 1:
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice[0],
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			default:
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    multi,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			}
		} else {
			switch len(valueSlice) {
			case 0:
				return nil, nil
			case 1:
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice[0],
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			default:
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_In,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			}
		}
	case reflect.Struct, reflect.Func, reflect.Map, reflect.Chan:
		return nil, gerror.Newf("Query field kind %s is not supported", t.Kind())
	case reflect.String:
		valueString := value.(string)
		if g.IsEmpty(valueString) {
			return nil, nil
		}
		if strings.HasPrefix(valueString, ConditionQueryPrefix) && strings.HasSuffix(valueString, ConditionQuerySuffix) {
			var condition *PropertyFilter
			err := gjson.DecodeTo(valueString[9:], &condition)
			condition.Property = columnName
			if err != nil {
				return nil, gwerror.NewBadRequestErrorf(req, err.Error())
			}
			g.Log().Debugf(ctx, "Query field type is orm.Condition: %s", gjson.MustEncodeString(condition))
			return condition, nil
		}
		wildcardString, ok := tag.Lookup(TagNameWildcard)
		if ok {
			wildcard, err := gwtypes.ParseWildcardType(wildcardString)
			if err != nil {
				return nil, err
			}
			switch wildcard {
			case gwtypes.WildcardType_Contains:
				return &PropertyFilter{
					Property: columnName,
					Value:    decorateValueStrForWildcard(valueString, gwtypes.WildcardType_Contains, mType),
					Operator: gwtypes.OperatorType_Like,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: wildcard,
				}, nil
			case gwtypes.WildcardType_StartsWith:
				return &PropertyFilter{
					Property: columnName,
					Value:    decorateValueStrForWildcard(valueString, gwtypes.WildcardType_StartsWith, mType),
					Operator: gwtypes.OperatorType_Like,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: wildcard,
				}, nil
			case gwtypes.WildcardType_EndsWith:
				return &PropertyFilter{
					Property: columnName,
					Value:    decorateValueStrForWildcard(valueString, gwtypes.WildcardType_EndsWith, mType),
					Operator: gwtypes.OperatorType_Like,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: wildcard,
				}, nil
			default:
				return &PropertyFilter{
					Property: columnName,
					Value:    valueString,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: wildcard,
				}, nil
			}
		} else {
			return &PropertyFilter{
				Property: columnName,
				Value:    valueString,
				Operator: gwtypes.OperatorType_EQ,
				Multi:    gwtypes.MultiType_Exact,
				Wildcard: gwtypes.WildcardType_None,
			}, nil
		}
	default:
		return &PropertyFilter{
			Property: columnName,
			Value:    value,
			Operator: gwtypes.OperatorType_EQ,
			Multi:    gwtypes.MultiType_Exact,
			Wildcard: gwtypes.WildcardType_None,
		}, nil
	}
}

func unwrapAnyFilter(columnName string, tag reflect.StructTag, valueAny *anypb.Any, mType ModelType) (pf *PropertyFilter, err error) {
	if valueAny == nil {
		return nil, nil
	}
	v, err := valueAny.UnmarshalNew()
	if err != nil {
		return nil, nil
	}

	switch vt := v.(type) {
	case *gwtypes.BoolSlice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *gwtypes.DoubleSlice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *gwtypes.FloatSlice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *gwtypes.UInt32Slice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *gwtypes.UInt64Slice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *gwtypes.Int32Slice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *gwtypes.Int64Slice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *gwtypes.StringSlice:
		return parseFieldSliceFilter(columnName, tag, vt.Value)
	case *wrapperspb.BoolValue:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.BytesValue:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.DoubleValue:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.FloatValue:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.Int32Value:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.Int64Value:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.UInt32Value:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.UInt64Value:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *wrapperspb.StringValue:
		return parseFieldSingleFilter(columnName, vt.Value)
	case *gwtypes.Condition:
		return AddConditionFilter(columnName, vt, mType)
	default:
		return nil, gerror.Newf("Unsupported value type: %v", vt)
	}
}

func parseFieldSingleFilter(columnName string, value interface{}) (pf *PropertyFilter, err error) {
	if value == nil {
		return nil, nil
	}
	return &PropertyFilter{
		Property: columnName,
		Value:    value,
		Operator: gwtypes.OperatorType_EQ,
		Multi:    gwtypes.MultiType_Exact,
		Wildcard: gwtypes.WildcardType_None,
	}, nil
}

func parseFieldSliceFilter[T any](columnName string, tag reflect.StructTag, value []T) (pf *PropertyFilter, err error) {
	if value == nil {
		return nil, nil
	}
	if multiTag, ok := tag.Lookup(TagNameMulti); ok {
		multi, err := gwtypes.ParseMultiType(multiTag)
		if err != nil {
			return nil, err
		}
		switch len(value) {
		case 0:
			return nil, nil
		case 1:
			return &PropertyFilter{
				Property: columnName,
				Value:    value[0],
				Operator: gwtypes.OperatorType_EQ,
				Multi:    gwtypes.MultiType_Exact,
				Wildcard: gwtypes.WildcardType_None,
			}, nil
		case 2:
			if multi == gwtypes.MultiType_Between {
				return &PropertyFilter{
					Property: columnName,
					Value:    value,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Between,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else {
				return &PropertyFilter{
					Property: columnName,
					Value:    value,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_In,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			}
		default:
			return &PropertyFilter{
				Property: columnName,
				Value:    value,
				Operator: gwtypes.OperatorType_EQ,
				Multi:    gwtypes.MultiType_In,
				Wildcard: gwtypes.WildcardType_None,
			}, nil
		}
	} else {
		switch len(value) {
		case 0:
			return nil, nil
		case 1:
			return &PropertyFilter{
				Property: columnName,
				Value:    value[0],
				Operator: gwtypes.OperatorType_EQ,
				Multi:    gwtypes.MultiType_Exact,
				Wildcard: gwtypes.WildcardType_None,
			}, nil
		default:
			return &PropertyFilter{
				Property: columnName,
				Value:    value,
				Operator: gwtypes.OperatorType_EQ,
				Multi:    gwtypes.MultiType_In,
				Wildcard: gwtypes.WildcardType_None,
			}, nil
		}
	}
}

func AddConditionFilter(columnName string, condition *gwtypes.Condition, mType ModelType) (pf *PropertyFilter, err error) {
	if condition == nil {
		return nil, nil
	}
	// todo 当 condition.Operator 为 Null、NotNull 时，允许 nil 的 Value
	if condition.Value == nil {
		return nil, nil
	}
	v, err := condition.Value.UnmarshalNew()
	if err != nil {
		return nil, nil
	}
	switch vt := v.(type) {
	case *gwtypes.BoolSlice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.BoolSlice).Value)
	case *gwtypes.DoubleSlice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.DoubleSlice).Value)
	case *gwtypes.FloatSlice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.FloatSlice).Value)
	case *gwtypes.UInt32Slice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.UInt32Slice).Value)
	case *gwtypes.UInt64Slice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.UInt64Slice).Value)
	case *gwtypes.Int32Slice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.Int32Slice).Value)
	case *gwtypes.Int64Slice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.Int64Slice).Value)
	case *gwtypes.StringSlice:
		return parseFieldConditionSliceFilter(columnName, condition, v.(*gwtypes.StringSlice).Value)
	case *wrapperspb.BoolValue:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.BoolValue).Value, mType)
	case *wrapperspb.BytesValue:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.BytesValue).Value, mType)
	case *wrapperspb.DoubleValue:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.DoubleValue).Value, mType)
	case *wrapperspb.FloatValue:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.FloatValue).Value, mType)
	case *wrapperspb.Int32Value:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.Int32Value).Value, mType)
	case *wrapperspb.Int64Value:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.Int64Value).Value, mType)
	case *wrapperspb.UInt32Value:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.UInt32Value).Value, mType)
	case *wrapperspb.UInt64Value:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.UInt64Value).Value, mType)
	case *wrapperspb.StringValue:
		return parseFieldConditionSingleFilter(columnName, condition, v.(*wrapperspb.StringValue).Value, mType)
	default:
		return nil, gerror.Newf("不支持的Value类型%v", vt)
	}
}

func parseFieldConditionSliceFilter[T any](columnName string, condition *gwtypes.Condition, valueSlice []T) (*PropertyFilter, error) {
	valueLen := len(valueSlice)
	if valueLen == 0 {
		return nil, nil
	}
	switch condition.Operator {
	case gwtypes.OperatorType_EQ:
		switch condition.Multi {
		case gwtypes.MultiType_Exact:
			if valueLen == 1 {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice[0],
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_In,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			}
		case gwtypes.MultiType_Between:
			if valueLen == 1 {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice[0],
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else if valueLen == 2 {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Between,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else {
				return nil, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", columnName, valueLen)
			}
		case gwtypes.MultiType_NotBetween:
			if valueLen == 1 {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice[0],
					Operator: gwtypes.OperatorType_NE,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else if valueLen == 2 {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_NotBetween,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else {
				return nil, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", columnName, valueLen)
			}
		case gwtypes.MultiType_In:
			if valueLen == 1 {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice[0],
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_In,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			}
		case gwtypes.MultiType_NotIn:
			if valueLen == 1 {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice[0],
					Operator: gwtypes.OperatorType_NE,
					Multi:    gwtypes.MultiType_Exact,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			} else {
				return &PropertyFilter{
					Property: columnName,
					Value:    valueSlice,
					Operator: gwtypes.OperatorType_EQ,
					Multi:    gwtypes.MultiType_NotIn,
					Wildcard: gwtypes.WildcardType_None,
				}, nil
			}
		}
	case gwtypes.OperatorType_NE:
	case gwtypes.OperatorType_GT:
	case gwtypes.OperatorType_GTE:
	case gwtypes.OperatorType_LT:
	case gwtypes.OperatorType_LTE:
	case gwtypes.OperatorType_Like:
	case gwtypes.OperatorType_NotLike:
		return nil, gerror.Newf("Operator值为'%s'，但传入的Value: '%s'不应该是数组", gwtypes.OperatorType_name[int32(condition.Operator)], gconv.String(valueSlice))
	case gwtypes.OperatorType_Null:
	case gwtypes.OperatorType_NotNull:
		return nil, gerror.Newf("Operator值为'%s'，但传入的Value: '%s'应该为nil", gwtypes.OperatorType_name[int32(condition.Operator)], gconv.String(valueSlice))
	default:
		return nil, gerror.Newf("不支持的Operator值，传入Value: '%s'", gwtypes.OperatorType_name[int32(condition.Operator)], gconv.String(valueSlice))
	}
	return nil, nil
}

func parseFieldConditionSingleFilter(columnName string, condition *gwtypes.Condition, value interface{}, mType ModelType) (*PropertyFilter, error) {
	if value == nil && condition.Operator != gwtypes.OperatorType_Null && condition.Operator != gwtypes.OperatorType_NotNull {
		return nil, nil
	}
	switch condition.Operator {
	case gwtypes.OperatorType_EQ:
		switch condition.Multi {
		case gwtypes.MultiType_Exact:
			return &PropertyFilter{
				Property: columnName,
				Value:    value,
				Operator: gwtypes.OperatorType_EQ,
				Multi:    gwtypes.MultiType_Exact,
				Wildcard: gwtypes.WildcardType_None,
			}, nil
		case gwtypes.MultiType_Between:
			return nil, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_Between)], gconv.String(value))
		case gwtypes.MultiType_NotBetween:
			return nil, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_NotBetween)], gconv.String(value))
		case gwtypes.MultiType_In:
			return nil, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_In)], gconv.String(value))
		case gwtypes.MultiType_NotIn:
			return nil, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_NotIn)], gconv.String(value))
		}
	case gwtypes.OperatorType_NE, gwtypes.OperatorType_GT, gwtypes.OperatorType_GTE, gwtypes.OperatorType_LT, gwtypes.OperatorType_LTE:
		return &PropertyFilter{
			Property: columnName,
			Value:    value,
			Operator: condition.Operator,
			Multi:    gwtypes.MultiType_Exact,
			Wildcard: gwtypes.WildcardType_None,
		}, nil
	case gwtypes.OperatorType_Like, gwtypes.OperatorType_NotLike:
		valueStr := gconv.String(value)
		if g.IsEmpty(valueStr) {
			return nil, nil
		}
		valueStr = decorateValueStrForWildcard(valueStr, condition.Wildcard, mType)
		return &PropertyFilter{
			Property: columnName,
			Value:    valueStr,
			Operator: condition.Operator,
			Multi:    gwtypes.MultiType_Exact,
			Wildcard: condition.Wildcard,
		}, nil
	case gwtypes.OperatorType_Null, gwtypes.OperatorType_NotNull:
		return &PropertyFilter{
			Property: columnName,
			Value:    nil,
			Operator: condition.Operator,
			Multi:    gwtypes.MultiType_Exact,
			Wildcard: gwtypes.WildcardType_None,
		}, nil
	}
	return nil, nil
}

// ParseConditions 根据传入 query 结构指针的值来设定 Where 条件
// ctx         context
// queryPtr    传入 query 结构指针的值
// columnMap   表的字段定义map，key为GoField，value为表字段名
// m           gdb.Model
func ParseConditions(ctx context.Context, req interface{}, columnMap map[string]string, m *Model) (*Model, error) {
	var err error
	p := reflect.TypeOf(req)
	if p.Kind() != reflect.Ptr { // 要求传入值必须是个指针
		return m, gwerror.NewBadRequestErrorf(req, "服务函数的输入参数必须是结构体指针")
	}
	t := p.Elem()
	//g.Log().Debugf(ctx, "kind of input parameter is %s", t.Name())

	queryValue := reflect.ValueOf(req).Elem()

	// 循环结构体的字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		if fieldName == "Meta" {
			continue
		}
		// Only do converting to public attributes.
		if !field.IsExported() {
			continue
		}
		fieldType := field.Type
		//g.Log().Debugf(ctx, "kind of field \"%s\" is %s", fieldName, field.Type.Kind().String())
		switch fieldType.Kind() {
		case reflect.Ptr:
			columnName, exists := columnMap[fieldName]
			if !exists {
				continue
			}
			fieldElemType := fieldType.Elem()
			//g.Log().Debugf(ctx, "kind of element of field %s is %s", fieldName, fieldElemType.Kind().String())
			if fieldElemType == reflect.TypeOf(anypb.Any{}) {
				anyValue := queryValue.Field(i).Interface().(*anypb.Any)
				if anyValue != nil {
					m, err = unwrapAny(columnName, field.Tag, anyValue, m)
					if err != nil {
						return nil, err
					}
				}
			}
		case reflect.Struct:
			structValue := queryValue.Field(i)
			g.Log().Debugf(ctx, "value of field %s is %x", fieldName, structValue)
			for si := 0; si < fieldType.NumField(); si++ {
				innerField := fieldType.Field(si)
				if innerField.Type.Kind() != reflect.Interface { // 仅处理类型为 interface{} 的字段
					continue
				}
				columnName, exists := columnMap[innerField.Name] // 仅处理在表字段定义中有的字段
				if !exists {
					continue
				}
				fieldValue := structValue.Field(si).Interface()
				if fieldValue == nil { // 不出来值为nil的字段
					continue
				}
				g.Log().Debugf(ctx, "inner field %s kind:%si, column:%s, value:%s", innerField.Name, innerField.Type.Kind().String(), columnName, fieldValue)
				m, err = parseField(ctx, req, columnName, innerField.Tag, fieldValue, m)
				if err != nil {
					return nil, err
				}
			}
		case reflect.Interface:
			columnName, exists := columnMap[fieldName]
			if !exists {
				continue
			}
			fieldValue := queryValue.Field(i).Interface()
			if fieldValue == nil {
				continue
			}
			m, err = parseField(ctx, req, columnName, field.Tag, fieldValue, m)
			if err != nil {
				return nil, err
			}
		}
	}
	return m, nil
}

func unwrapAny(columnName string, tag reflect.StructTag, valueAny *anypb.Any, m *Model) (*Model, error) {
	if valueAny == nil {
		return m, nil
	}
	v, err := valueAny.UnmarshalNew()
	if err != nil {
		return m, nil
	}

	switch vt := v.(type) {
	case *gwtypes.BoolSlice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *gwtypes.DoubleSlice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *gwtypes.FloatSlice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *gwtypes.UInt32Slice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *gwtypes.UInt64Slice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *gwtypes.Int32Slice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *gwtypes.Int64Slice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *gwtypes.StringSlice:
		return parseFieldSlice(columnName, tag, vt.Value, m)
	case *wrapperspb.BoolValue:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.BytesValue:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.DoubleValue:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.FloatValue:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.Int32Value:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.Int64Value:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.UInt32Value:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.UInt64Value:
		return m.Where(columnName, vt.Value), nil
	case *wrapperspb.StringValue:
		return m.Where(columnName, vt.Value), nil
	case *gwtypes.Condition:
		return AddCondition1(columnName, vt, m)
	default:
		return nil, gerror.Newf("Unsupported value type: %v", vt)
	}
}

func parseField(ctx context.Context, req interface{}, columnName string, tag reflect.StructTag, value interface{}, m *Model) (*Model, error) {
	if value == nil {
		return m, nil
	}
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem() == reflect.TypeOf(PropertyFilter{}) {
			return AddCondition(ctx, columnName, value.(*PropertyFilter), m)
		}
		return m.Where(columnName, value), nil
	case reflect.Slice, reflect.Array:
		valueSlice := gconv.SliceAny(value)
		multiTag, ok := tag.Lookup(TagNameMulti)
		if ok {
			multi, err := gwtypes.ParseMultiType(multiTag)
			if err != nil {
				return nil, err
			}
			switch len(valueSlice) {
			case 0:
				return m, nil
			case 1:
				return m.Where(columnName, valueSlice[0]), nil
			case 2:
				if multi == gwtypes.MultiType_Between {
					return m.WhereBetween(columnName, valueSlice[0], valueSlice[1]), nil
				}
			default:
				return m.WhereIn(columnName, valueSlice), nil
			}
		} else {
			switch len(valueSlice) {
			case 0:
				return m, nil
			case 1:
				return m.Where(columnName, valueSlice[0]), nil
			default:
				return m.WhereIn(columnName, valueSlice), nil
			}
		}
	case reflect.Struct, reflect.Func, reflect.Map, reflect.Chan:
		g.Log().Warningf(ctx, "Query field type struct is not supported")
		return m, nil
	case reflect.String:
		valueString := value.(string)
		if g.IsEmpty(valueString) {
			return m, nil
		}
		if strings.HasPrefix(valueString, ConditionQueryPrefix) && strings.HasSuffix(valueString, ConditionQuerySuffix) {
			var condition *PropertyFilter
			err := gjson.DecodeTo(valueString[9:], &condition)
			if err != nil {
				return nil, gwerror.NewBadRequestErrorf(req, err.Error())
			} else {
				g.Log().Debugf(ctx, "Query field type is orm.PropertyFilter: %s", gjson.MustEncodeString(condition))
				return AddCondition(ctx, columnName, condition, m)
			}
		}
		wildcardString, ok := tag.Lookup(TagNameWildcard)
		if ok {
			wildcard, err := gwtypes.ParseWildcardType(wildcardString)
			if err != nil {
				return nil, err
			}
			valueString = decorateValueStrForWildcard(valueString, wildcard, m.Type)
			switch wildcard {
			case gwtypes.WildcardType_Contains:
			case gwtypes.WildcardType_StartsWith:
			case gwtypes.WildcardType_EndsWith:
				return m.WhereLike(columnName, valueString), nil
			default:
				return m.WhereIn(columnName, valueString), nil
			}
		} else {
			return m.Where(columnName, valueString), nil
		}
	default:
		return m.Where(columnName, value), nil
	}
	return m, nil
}

func parseFieldSlice[T any](columnName string, tag reflect.StructTag, value []T, m *Model) (*Model, error) {
	if value == nil {
		return m, nil
	}
	if multiTag, ok := tag.Lookup(TagNameMulti); ok {
		multi, err := gwtypes.ParseMultiType(multiTag)
		if err != nil {
			return nil, err
		}
		switch len(value) {
		case 0:
			return m, nil
		case 1:
			return m.Where(columnName, value[0]), nil
		case 2:
			if multi == gwtypes.MultiType_Between {
				return m.WhereBetween(columnName, value[0], value[1]), nil
			}
		default:
			return m.WhereIn(columnName, value), nil
		}
	} else {
		switch len(value) {
		case 0:
			return m, nil
		case 1:
			return m.Where(columnName, value[0]), nil
		default:
			return m.WhereIn(columnName, value), nil
		}
	}
	return m, nil
}

func parseFieldConditionSingle(columnName string, condition *gwtypes.Condition, value interface{}, m *Model) (*Model, error) {
	if value == nil && condition.Operator != gwtypes.OperatorType_Null && condition.Operator != gwtypes.OperatorType_NotNull {
		return m, nil
	}
	switch condition.Operator {
	case gwtypes.OperatorType_EQ:
		switch condition.Multi {
		case gwtypes.MultiType_Exact:
			return m.WhereIn(columnName, value), nil
		case gwtypes.MultiType_Between:
			return m, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_Between)], gconv.String(value))
		case gwtypes.MultiType_NotBetween:
			return m, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_NotBetween)], gconv.String(value))
		case gwtypes.MultiType_In:
			return m, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_In)], gconv.String(value))
		case gwtypes.MultiType_NotIn:
			return m, gerror.Newf("Multi值为'%s'，但传入的Value: '%s'并非数组", gwtypes.MultiType_name[int32(gwtypes.MultiType_NotIn)], gconv.String(value))
		}
	case gwtypes.OperatorType_NE:
		return m.WhereNot(columnName, value), nil
	case gwtypes.OperatorType_GT:
		return m.WhereGT(columnName, value), nil
	case gwtypes.OperatorType_GTE:
		return m.WhereGTE(columnName, value), nil
	case gwtypes.OperatorType_LT:
		return m.WhereLT(columnName, value), nil
	case gwtypes.OperatorType_LTE:
		return m.WhereLTE(columnName, value), nil
	case gwtypes.OperatorType_Like:
		valueStr := gconv.String(value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		valueStr = decorateValueStrForWildcard(valueStr, condition.Wildcard, m.Type)
		return m.WhereLike(columnName, valueStr), nil
	case gwtypes.OperatorType_NotLike:
		valueStr := gconv.String(value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		valueStr = decorateValueStrForWildcard(valueStr, condition.Wildcard, m.Type)
		return m.WhereNotLike(columnName, valueStr), nil
	case gwtypes.OperatorType_Null:
		return m.WhereNull(columnName), nil
	case gwtypes.OperatorType_NotNull:
		return m.WhereNotNull(columnName), nil
	}
	return m, nil
}

func parseFieldConditionSlice[T any](columnName string, condition *gwtypes.Condition, valueSlice []T, m *Model) (*Model, error) {
	valueLen := len(valueSlice)
	if valueLen == 0 {
		return m, nil
	}
	switch condition.Operator {
	case gwtypes.OperatorType_EQ:
		switch condition.Multi {
		case gwtypes.MultiType_Exact:
			if valueLen == 1 {
				return m.Where(columnName, valueSlice[0]), nil
			} else {
				return m.WhereIn(columnName, valueSlice), nil
			}
		case gwtypes.MultiType_Between:
			if valueLen == 1 {
				return m.Where(columnName, valueSlice[0]), nil
			} else if valueLen == 2 {
				return m.WhereBetween(columnName, valueSlice[0], valueSlice[1]), nil
			} else {
				return m, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", columnName, valueLen)
			}
		case gwtypes.MultiType_NotBetween:
			if valueLen == 1 {
				return m.WhereNot(columnName, valueSlice[0]), nil
			} else if valueLen == 2 {
				return m.WhereNotBetween(columnName, valueSlice[0], valueSlice[1]), nil
			} else {
				return m, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", columnName, valueLen)
			}
		case gwtypes.MultiType_In:
			if valueLen == 1 {
				return m.Where(columnName, valueSlice[0]), nil
			} else {
				return m.WhereIn(columnName, valueSlice), nil
			}
		case gwtypes.MultiType_NotIn:
			if valueLen == 1 {
				return m.WhereNot(columnName, valueSlice[0]), nil
			} else {
				return m.WhereNotIn(columnName, valueSlice), nil
			}
		}
	case gwtypes.OperatorType_NE:
	case gwtypes.OperatorType_GT:
	case gwtypes.OperatorType_GTE:
	case gwtypes.OperatorType_LT:
	case gwtypes.OperatorType_LTE:
	case gwtypes.OperatorType_Like:
	case gwtypes.OperatorType_NotLike:
		return m, gerror.Newf("Operator值为'%s'，但传入的Value: '%s'不应该是数组", gwtypes.OperatorType_name[int32(condition.Operator)], gconv.String(valueSlice))
	case gwtypes.OperatorType_Null:
	case gwtypes.OperatorType_NotNull:
		return m, gerror.Newf("Operator值为'%s'，但传入的Value: '%s'应该为nil", gwtypes.OperatorType_name[int32(condition.Operator)], gconv.String(valueSlice))
	default:
		return m, gerror.Newf("不支持的Operator值，传入Value: '%s'", gwtypes.OperatorType_name[int32(condition.Operator)], gconv.String(valueSlice))
	}
	return m, nil
}

func AddCondition1(columnName string, condition *gwtypes.Condition, m *Model) (*Model, error) {
	if condition == nil {
		return m, nil
	}
	// todo 当 condition.Operator 为 Null、NotNull 时，允许 nil 的 Value
	if condition.Value == nil {
		return m, nil
	}
	v, err := condition.Value.UnmarshalNew()
	if err != nil {
		return m, nil
	}
	switch vt := v.(type) {
	case *gwtypes.BoolSlice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.BoolSlice).Value, m)
	case *gwtypes.DoubleSlice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.DoubleSlice).Value, m)
	case *gwtypes.FloatSlice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.FloatSlice).Value, m)
	case *gwtypes.UInt32Slice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.UInt32Slice).Value, m)
	case *gwtypes.UInt64Slice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.UInt64Slice).Value, m)
	case *gwtypes.Int32Slice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.Int32Slice).Value, m)
	case *gwtypes.Int64Slice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.Int64Slice).Value, m)
	case *gwtypes.StringSlice:
		return parseFieldConditionSlice(columnName, condition, v.(*gwtypes.StringSlice).Value, m)
	case *wrapperspb.BoolValue:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.BoolValue).Value, m)
	case *wrapperspb.BytesValue:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.BytesValue).Value, m)
	case *wrapperspb.DoubleValue:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.DoubleValue).Value, m)
	case *wrapperspb.FloatValue:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.FloatValue).Value, m)
	case *wrapperspb.Int32Value:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.Int32Value).Value, m)
	case *wrapperspb.Int64Value:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.Int64Value).Value, m)
	case *wrapperspb.UInt32Value:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.UInt32Value).Value, m)
	case *wrapperspb.UInt64Value:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.UInt64Value).Value, m)
	case *wrapperspb.StringValue:
		return parseFieldConditionSingle(columnName, condition, v.(*wrapperspb.StringValue).Value, m)
	default:
		return m, gerror.Newf("不支持的Value类型%v", vt)
	}
}

func AddCondition(_ context.Context, columnName string, condition *PropertyFilter, m *Model) (*Model, error) {
	if condition == nil {
		return m, nil
	}
	if condition.Value == nil {
		return m, nil
	}
	switch condition.Operator {
	case gwtypes.OperatorType_EQ:
		switch condition.Multi {
		case gwtypes.MultiType_Exact:
			return m.Where(columnName, condition.Value), nil
		case gwtypes.MultiType_Between:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.Where(columnName, valueSlice[0]), nil
			} else if valueLen == 2 {
				return m.WhereBetween(columnName, valueSlice[0], valueSlice[1]), nil
			} else {
				return m, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", columnName, valueLen)
			}
		case gwtypes.MultiType_NotBetween:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.WhereNot(columnName, valueSlice[0]), nil
			} else if valueLen == 2 {
				return m.WhereNotBetween(columnName, valueSlice[0], valueSlice[1]), nil
			} else {
				return m, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", columnName, valueLen)
			}
		case gwtypes.MultiType_In:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.Where(columnName, valueSlice[0]), nil
			} else {
				return m.WhereIn(columnName, valueSlice), nil
			}
		case gwtypes.MultiType_NotIn:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.WhereNot(columnName, valueSlice[0]), nil
			} else {
				return m.WhereNotIn(columnName, valueSlice), nil
			}
		}
	case gwtypes.OperatorType_NE:
		return m.WhereNot(columnName, condition.Value), nil
	case gwtypes.OperatorType_GT:
		return m.WhereGT(columnName, condition.Value), nil
	case gwtypes.OperatorType_GTE:
		return m.WhereGTE(columnName, condition.Value), nil
	case gwtypes.OperatorType_LT:
		return m.WhereLT(columnName, condition.Value), nil
	case gwtypes.OperatorType_LTE:
		return m.WhereLTE(columnName, condition.Value), nil
	case gwtypes.OperatorType_Like:
		valueStr := gconv.String(condition.Value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		valueStr = decorateValueStrForWildcard(valueStr, condition.Wildcard, m.Type)
		return m.WhereLike(columnName, valueStr), nil
	case gwtypes.OperatorType_NotLike:
		valueStr := gconv.String(condition.Value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		valueStr = decorateValueStrForWildcard(valueStr, condition.Wildcard, m.Type)
		return m.WhereNotLike(columnName, valueStr), nil
	case gwtypes.OperatorType_Null:
		return m.WhereNot(columnName, condition.Value), nil
	case gwtypes.OperatorType_NotNull:
		return m.WhereNot(columnName, condition.Value), nil
	}
	return m, nil
}

func decorateValueStrForWildcard(valueStr string, wildcardType gwtypes.WildcardType, modelType ModelType) string {
	switch wildcardType {
	case gwtypes.WildcardType_Contains:
		if modelType == GF_ORM {
			return "%" + valueStr + "%"
		} else if modelType == MONGO {
			return valueStr
		}
	case gwtypes.WildcardType_StartsWith:
		if modelType == GF_ORM {
			return valueStr + "%"
		} else if modelType == MONGO {
			return "^" + valueStr
		}
	case gwtypes.WildcardType_EndsWith:
		if modelType == GF_ORM {
			return "%" + valueStr
		} else if modelType == MONGO {
			return valueStr + "$"
		}
	}
	return valueStr
}

func ApplyFilter(ctx context.Context, request FilterRequest, m *Model) (output *Model, err error) {
	output = m
	pfs := request.PropertyFilters
	if len(pfs) == 0 {
		return
	}
	for _, pf := range pfs {
		output, err = ApplyPropertyFilter(ctx, pf, m)
	}
	return
}

func ApplyPropertyFilter(_ context.Context, condition *PropertyFilter, m *Model) (*Model, error) {
	if condition == nil {
		return m, nil
	}
	if condition.Value == nil {
		return m, nil
	}
	property := condition.Property
	switch condition.Operator {
	case gwtypes.OperatorType_EQ:
		switch condition.Multi {
		case gwtypes.MultiType_Exact:
			return m.Where(property, condition.Value), nil
		case gwtypes.MultiType_Between:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.Where(property, valueSlice[0]), nil
			} else if valueLen == 2 {
				return m.WhereBetween(property, valueSlice[0], valueSlice[1]), nil
			} else {
				return m, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", property, valueLen)
			}
		case gwtypes.MultiType_NotBetween:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.WhereNot(property, valueSlice[0]), nil
			} else if valueLen == 2 {
				return m.WhereNotBetween(property, valueSlice[0], valueSlice[1]), nil
			} else {
				return m, gwerror.NewBadRequestErrorf("column %s requires between query but given %d values", property, valueLen)
			}
		case gwtypes.MultiType_In:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.Where(property, valueSlice[0]), nil
			} else {
				return m.WhereIn(property, valueSlice), nil
			}
		case gwtypes.MultiType_NotIn:
			valueSlice := gconv.SliceAny(condition.Value)
			valueLen := len(valueSlice)
			if valueLen == 0 {
				return m, nil
			} else if valueLen == 1 {
				return m.WhereNot(property, valueSlice[0]), nil
			} else {
				return m.WhereNotIn(property, valueSlice), nil
			}
		}
	case gwtypes.OperatorType_NE:
		return m.WhereNot(property, condition.Value), nil
	case gwtypes.OperatorType_GT:
		return m.WhereGT(property, condition.Value), nil
	case gwtypes.OperatorType_GTE:
		return m.WhereGTE(property, condition.Value), nil
	case gwtypes.OperatorType_LT:
		return m.WhereLT(property, condition.Value), nil
	case gwtypes.OperatorType_LTE:
		return m.WhereLTE(property, condition.Value), nil
	case gwtypes.OperatorType_Like:
		valueStr := gconv.String(condition.Value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		valueStr = decorateValueStrForWildcard(valueStr, condition.Wildcard, m.Type)
		return m.WhereLike(property, valueStr), nil
	case gwtypes.OperatorType_NotLike:
		valueStr := gconv.String(condition.Value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		valueStr = decorateValueStrForWildcard(valueStr, condition.Wildcard, m.Type)
		return m.WhereNotLike(property, valueStr), nil
	case gwtypes.OperatorType_Null:
		return m.WhereNot(property, condition.Value), nil
	case gwtypes.OperatorType_NotNull:
		return m.WhereNot(property, condition.Value), nil
	}
	return m, nil
}
