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

package trace

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/WesleyWu/gowing/rpc/dubbogo/keys"
	"github.com/WesleyWu/gowing/util/gwmap"
	"github.com/gogf/gf/v2/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	extension.SetFilter(keys.SyncTraceFilterKey, func() filter.Filter {
		return &serverTraceFilter{}
	})
}

type serverTraceFilter struct {
}

func (f *serverTraceFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	traceId := gtrace.GetTraceID(ctx)
	if traceId == "" {
		attachments := invocation.Attachments()
		ctx = otel.GetTextMapPropagator().Extract(ctx, &gwmap.StrStrMap{
			InnerMap: attachments,
		})
		spanCtx := trace.SpanContextFromContext(ctx)
		traceId = spanCtx.TraceID().String()
		if traceId != "" {
			ctx, _ = gtrace.WithTraceID(ctx, traceId)
		}
	}
	return invoker.Invoke(ctx, invocation)
}
func (f *serverTraceFilter) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	return result
}
