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

package validation

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/WesleyWu/gowing/rpc/dubbogo/keys"
	"github.com/WesleyWu/gowing/util/gwreflect"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
)

const validationTag = "validation"

func init() {
	extension.SetFilter(keys.ValidationFilterKey, newFilter)
}

type validationFilter struct {
}

func newFilter() filter.Filter {
	return &validationFilter{}
}

func validateRequest(ctx context.Context, req interface{}) error {
	ctx, span := gtrace.NewSpan(ctx, "validateRequest")
	defer span.End()
	return g.Validator().Data(req).Run(ctx)
}

func shouldValidate(ctx context.Context, param interface{}) bool {
	metaField, err := gwreflect.GetMetaField(ctx, param)
	if err != nil {
		return false
	}
	if metaField == nil {
		return false
	}
	return metaField.Tag.Get(validationTag) == "true"
}

func (f *validationFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	ctx, span := gtrace.NewSpan(ctx, "validationFilter.Invoke")
	defer span.End()
	params := invocation.Arguments()
	for _, param := range params {
		if !shouldValidate(ctx, param) {
			continue
		}
		err := validateRequest(ctx, param)
		if err != nil {
			return &protocol.RPCResult{
				Attrs: nil,
				Err:   err,
				Rest:  err.Error(),
			}
		}
		//g.Log().Infof(ctx, "param %d: %v", i, param)
	}
	return invoker.Invoke(ctx, invocation)
}
func (f *validationFilter) OnResponse(ctx context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	err := result.Error()
	if err != nil {
		g.Log().Infof(ctx, "validationFilter OnResponse error: %v", result.Error())
	}
	return result
}
