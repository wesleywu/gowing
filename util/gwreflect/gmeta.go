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

package gwreflect

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gstructs"
	"reflect"
)

const (
	metaAttributeName = "Meta"
	metaTypeName      = "gmeta.Meta" // metaTypeName is for type string comparison.
)

func GetMetaField(ctx context.Context, req interface{}) (*reflect.StructField, error) {
	ctx, span := gtrace.NewSpan(ctx, "GetMetaField")
	defer span.End()
	reflectType, err := gstructs.StructType(req)
	if err != nil {
		return nil, err
	}
	metaField, ok := reflectType.FieldByName(metaAttributeName)
	if !ok {
		return nil, nil
	}
	if metaField.Type.String() != metaTypeName {
		return nil, nil
	}
	return &metaField, nil
}
