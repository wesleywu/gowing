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
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
	"reflect"
)

func MergeDefaultStructValue(_ context.Context, pointer interface{}) error {
	defaultValueTags := []string{gtag.DefaultShort, gtag.Default}
	tagFields, err := gstructs.TagFields(pointer, defaultValueTags)
	if err != nil {
		return err
	}
	if len(tagFields) > 0 {
		for _, field := range tagFields {
			if field.Value.IsZero() {
				//fieldValue := reflect.ValueOf(field.TagValue)
				tagValueConverted := gconv.Convert(field.TagValue, field.Type().String())
				field.Value.Set(reflect.ValueOf(tagValueConverted))
			}
		}
	}
	return nil
}
