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

package gwconv

import (
	"context"
	"github.com/WesleyWu/gowing/util/gwwrapper"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
)

var anyType = reflect.TypeOf(anypb.Any{})

func ToProtoMessage[T proto.Message](ctx context.Context, ptr interface{}, allowAllNilFields bool) (result T, err error) {
	var (
		srcType        reflect.Type
		srcElemType    reflect.Type
		srcElemValue   reflect.Value
		srcField       reflect.StructField
		srcFieldValue  reflect.Value
		destType       reflect.Type
		destValue      reflect.Value
		destField      reflect.StructField
		destFieldValue reflect.Value
		ok             bool
		allFieldsNil   = true
	)
	srcType = reflect.TypeOf(ptr)
	if srcType.Kind() != reflect.Ptr {
		g.Log().Warningf(ctx, "Parameter ptr should be a pointer of struct, but %s", srcType)
		return
	}
	srcElemType = srcType.Elem()
	if srcElemType.Kind() != reflect.Struct {
		g.Log().Warningf(ctx, "Parameter ptr should be a pointer of struct, but a pointer of %s", srcElemType)
		return
	}

	destType = reflect.TypeOf(result).Elem()
	destValue = reflect.New(destType)
	result, ok = destValue.Interface().(T)
	if !ok {
		g.Log().Warningf(ctx, "Output parameter should implement proto.Message, but was type '%s'", destType)
		return
	}
	destFieldMap := make(map[string]reflect.StructField)
	destFieldValueMap := make(map[string]reflect.Value)
	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		fieldName := field.Name
		if !field.IsExported() {
			continue
		}
		if fieldName == "Meta" {
			continue
		}
		destFieldMap[fieldName] = field
		destFieldValueMap[fieldName] = destValue.Elem().Field(i)
	}

	srcElemValue = reflect.ValueOf(ptr).Elem()
	for i := 0; i < srcElemType.NumField(); i++ {
		srcField = srcElemType.Field(i)
		fieldName := srcField.Name
		if !srcField.IsExported() {
			continue
		}
		if fieldName == "Meta" {
			continue
		}
		destField, ok = destFieldMap[fieldName]
		if !ok {
			continue
		}
		destFieldValue, ok = destFieldValueMap[fieldName]
		if !ok {
			continue
		}
		if destField.Type.Kind() == reflect.Ptr && destField.Type.Elem() == anyType {
			valueAny := gwwrapper.AnyInterfaceCondition(ctx, srcElemValue.Field(i).Interface(), destField.Tag)
			if valueAny != nil {
				destFieldValue.Set(reflect.ValueOf(valueAny))
				allFieldsNil = false
			}
		} else {
			srcFieldValue = srcElemValue.Field(i)
			if srcFieldValue.IsZero() {
				srcFieldValue = defaultValueIfZero(srcField, srcFieldValue)
			} else {
				allFieldsNil = false
			}
			if srcField.Type.Kind() == reflect.Interface { // 源字段为 interface 时，需要做转换
				srcFieldValue = convertToType(srcFieldValue, destField.Type)
			}
			destFieldValue.Set(srcFieldValue)
		}
	}
	if !allowAllNilFields && allFieldsNil {
		err = gerror.Newf("all field is nil")
	}
	return
}

func defaultValueIfZero(field reflect.StructField, value reflect.Value) reflect.Value {
	if !value.IsZero() {
		return value
	}
	defaultTag := field.Tag.Get(gtag.DefaultShort)
	if defaultTag == "" {
		defaultTag = field.Tag.Get(gtag.Default)
	}
	if defaultTag == "" {
		return value
	}
	tagValueConverted := gconv.Convert(defaultTag, field.Type.String())
	return reflect.ValueOf(tagValueConverted)
}

func convertToType(value reflect.Value, fieldType reflect.Type) reflect.Value {
	converted := gconv.Convert(value.Interface(), fieldType.String())
	return reflect.ValueOf(converted)
}

func FromProtoMessage[T proto.Message, Out interface{}](ctx context.Context, message T, result Out) error {
	err := gconv.Struct(message, result)
	if err != nil {
		return err
	}
	return nil
}
