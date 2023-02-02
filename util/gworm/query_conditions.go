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
	"github.com/WesleyWu/gowing/errors/gwerror"
	"github.com/WesleyWu/gowing/protobuf/gen/proto/go/gwtypes/v1"
	"github.com/WesleyWu/gowing/util/gwstring"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
	"strings"
)

const (
	ConditionQueryPrefix = "condition{"
	ConditionQuerySuffix = "}"
	TagNameMulti         = "multi"
	TagNameWildcard      = "wildcard"
)

type Condition struct {
	Operator gwtypes.OperatorType `json:"operator"`
	Multi    gwtypes.MultiType    `json:"multi"`
	Wildcard gwtypes.WildcardType `json:"wildcard"`
	Value    interface{}          `json:"value"`
}

// ParseConditions 根据传入 query 结构指针的值来设定 Where 条件
// ctx         context
// queryPtr    传入 query 结构指针的值
// columnMap   表的字段定义map，key为GoField，value为表字段名
// m           gdb.Model
func ParseConditions(ctx context.Context, req interface{}, columnMap map[string]string, m *gdb.Model) (*gdb.Model, error) {
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
		if !gwstring.IsLetterUpper(fieldName[0]) {
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

func unwrapAny(columnName string, tag reflect.StructTag, valueAny *anypb.Any, m *gdb.Model) (*gdb.Model, error) {
	if valueAny == nil {
		return m, nil
	}
	v, err := valueAny.UnmarshalNew()
	if err != nil {
		return m, nil
	}

	switch vt := v.(type) {
	case *gwtypes.BoolSlice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.BoolSlice).Value, m)
	case *gwtypes.DoubleSlice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.DoubleSlice).Value, m)
	case *gwtypes.FloatSlice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.FloatSlice).Value, m)
	case *gwtypes.UInt32Slice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.UInt32Slice).Value, m)
	case *gwtypes.UInt64Slice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.UInt64Slice).Value, m)
	case *gwtypes.Int32Slice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.Int32Slice).Value, m)
	case *gwtypes.Int64Slice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.Int64Slice).Value, m)
	case *gwtypes.StringSlice:
		return parseFieldSlice(columnName, tag, v.(*gwtypes.StringSlice).Value, m)
	case *wrapperspb.BoolValue:
		return m.Where(columnName, v.(*wrapperspb.BoolValue).Value), nil
	case *wrapperspb.BytesValue:
		return m.Where(columnName, v.(*wrapperspb.BytesValue).Value), nil
	case *wrapperspb.DoubleValue:
		return m.Where(columnName, v.(*wrapperspb.DoubleValue).Value), nil
	case *wrapperspb.FloatValue:
		return m.Where(columnName, v.(*wrapperspb.FloatValue).Value), nil
	case *wrapperspb.Int32Value:
		return m.Where(columnName, v.(*wrapperspb.Int32Value).Value), nil
	case *wrapperspb.Int64Value:
		return m.Where(columnName, v.(*wrapperspb.Int64Value).Value), nil
	case *wrapperspb.UInt32Value:
		return m.Where(columnName, v.(*wrapperspb.UInt32Value).Value), nil
	case *wrapperspb.UInt64Value:
		return m.Where(columnName, v.(*wrapperspb.UInt64Value).Value), nil
	case *wrapperspb.StringValue:
		return m.Where(columnName, v.(*wrapperspb.StringValue).Value), nil
	case *gwtypes.Condition:
		condition := v.(*gwtypes.Condition)
		return AddCondition1(columnName, condition, m)
	default:
		return nil, gerror.Newf("Unsupported value type: %v", vt)
	}
}

func parseField(ctx context.Context, req interface{}, columnName string, tag reflect.StructTag, value interface{}, m *gdb.Model) (*gdb.Model, error) {
	if value == nil {
		return m, nil
	}
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem() == reflect.TypeOf(Condition{}) {
			return AddCondition(ctx, columnName, value.(*Condition), m)
		}
		return m.Where(columnName, value), nil
	case reflect.Slice, reflect.Array:
		valueSlice := gconv.SliceAny(value)
		multiTag, ok := tag.Lookup(TagNameMulti)
		if ok {
			multi, err := ParseMultiType(multiTag)
			if err != nil {
				return nil, err
			}
			switch len(valueSlice) {
			case 0:
				return m, nil
			case 1:
				return m.Where(columnName, valueSlice[0]), nil
			case 2:
				if multi == Between {
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
			var condition *Condition
			err := gjson.DecodeTo(valueString[9:], &condition)
			if err != nil {
				return nil, gwerror.NewBadRequestErrorf(req, err.Error())
			} else {
				g.Log().Debugf(ctx, "Query field type is orm.Condition: %s", gjson.MustEncodeString(condition))
				return AddCondition(ctx, columnName, condition, m)
			}
		}
		wildcardString, ok := tag.Lookup(TagNameWildcard)
		if ok {
			wildcard, err := ParseWildcardType(wildcardString)
			if err != nil {
				return nil, err
			}
			switch wildcard {
			case Contains:
				return m.WhereLike(columnName, "%"+valueString+"%"), nil
			case StartsWith:
				return m.WhereLike(columnName, valueString+"%"), nil
			case EndsWith:
				return m.WhereLike(columnName, "%"+valueString), nil
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

func parseFieldSlice[T any](columnName string, tag reflect.StructTag, value []T, m *gdb.Model) (*gdb.Model, error) {
	if value == nil {
		return m, nil
	}
	if multiTag, ok := tag.Lookup(TagNameMulti); ok {
		multi, err := ParseMultiType(multiTag)
		if err != nil {
			return nil, err
		}
		switch len(value) {
		case 0:
			return m, nil
		case 1:
			return m.Where(columnName, value[0]), nil
		case 2:
			if multi == Between {
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

func parseFieldConditionSingle(columnName string, condition *gwtypes.Condition, value interface{}, m *gdb.Model) (*gdb.Model, error) {
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
		switch condition.Wildcard {
		case gwtypes.WildcardType_Contains:
			valueStr = "%" + valueStr + "%"
		case gwtypes.WildcardType_StartsWith:
			valueStr = valueStr + "%"
		case gwtypes.WildcardType_EndsWith:
			valueStr = "%" + valueStr
		}
		return m.WhereLike(columnName, valueStr), nil
	case gwtypes.OperatorType_NotLike:
		valueStr := gconv.String(value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		switch condition.Wildcard {
		case gwtypes.WildcardType_Contains:
			valueStr = "%" + valueStr + "%"
		case gwtypes.WildcardType_StartsWith:
			valueStr = valueStr + "%"
		case gwtypes.WildcardType_EndsWith:
			valueStr = "%" + valueStr
		}
		return m.WhereNotLike(columnName, valueStr), nil
	case gwtypes.OperatorType_Null:
		return m.WhereNull(columnName), nil
	case gwtypes.OperatorType_NotNull:
		return m.WhereNotNull(columnName), nil
	}
	return m, nil
}

func parseFieldConditionSlice[T any](columnName string, condition *gwtypes.Condition, valueSlice []T, m *gdb.Model) (*gdb.Model, error) {
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

func AddCondition1(columnName string, condition *gwtypes.Condition, m *gdb.Model) (*gdb.Model, error) {
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

func AddCondition(_ context.Context, columnName string, condition *Condition, m *gdb.Model) (*gdb.Model, error) {
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
		switch condition.Wildcard {
		case gwtypes.WildcardType_Contains:
			valueStr = "%" + valueStr + "%"
		case gwtypes.WildcardType_StartsWith:
			valueStr = valueStr + "%"
		case gwtypes.WildcardType_EndsWith:
			valueStr = "%" + valueStr
		}
		return m.WhereLike(columnName, valueStr), nil
	case gwtypes.OperatorType_NotLike:
		valueStr := gconv.String(condition.Value)
		if g.IsEmpty(valueStr) {
			return m, nil
		}
		switch condition.Wildcard {
		case gwtypes.WildcardType_Contains:
			valueStr = "%" + valueStr + "%"
		case gwtypes.WildcardType_StartsWith:
			valueStr = valueStr + "%"
		case gwtypes.WildcardType_EndsWith:
			valueStr = "%" + valueStr
		}
		return m.WhereNotLike(columnName, valueStr), nil
	case gwtypes.OperatorType_Null:
		return m.WhereNot(columnName, condition.Value), nil
	case gwtypes.OperatorType_NotNull:
		return m.WhereNot(columnName, condition.Value), nil
	}
	return m, nil
}
