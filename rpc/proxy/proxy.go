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

package proxy

import (
	"context"
	"reflect"

	"github.com/wesleywu/gowing/rpc/inproc/interceptor/cache"
	"github.com/wesleywu/gowing/rpc/inproc/types"
	"github.com/wesleywu/gowing/util/gwconv"
	"google.golang.org/protobuf/proto"
)

// NewInvocationProxy is a proxy function which wrap the original function with interceptors
// todo customize interceptors
func NewInvocationProxy[TReq proto.Message, TRes proto.Message](serviceName, methodName string, fn types.MethodFunc[TReq, TRes]) types.MethodFunc[TReq, TRes] {
	invocation := &types.MethodInvocation[TReq, TRes]{
		Method:      fn,
		ServiceName: serviceName,
		MethodName:  methodName,
	}
	invocation.Interceptors = append(invocation.Interceptors, cache.NewCacheInterceptor[TReq, TRes]().InterceptMethod)
	for _, interceptor := range invocation.Interceptors {
		invocation = interceptor(invocation)
	}
	return invocation.Method
}

func CallServiceMethod[TReq proto.Message, TRes proto.Message, Req interface{}, Res interface{}](
	ctx context.Context, req Req, allowAllNilFields bool, fn types.MethodFunc[TReq, TRes]) (res Res, err error) {
	var (
		in  TReq
		out TRes
	)
	in, err = gwconv.ToProtoMessage[TReq](ctx, req, allowAllNilFields)
	if err != nil {
		return
	}
	out, err = fn(ctx, in)
	if err != nil {
		return
	}
	res = reflect.New(reflect.TypeOf(res).Elem()).Interface().(Res)
	err = gwconv.FromProtoMessage[TRes](ctx, out, res)
	return
}
