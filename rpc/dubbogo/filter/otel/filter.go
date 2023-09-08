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

package otel

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	dubbogoTrace "dubbo.apache.org/dubbo-go/v3/filter/otel/trace"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/wesleywu/gowing/rpc/dubbogo/keys"
	"github.com/wesleywu/gowing/util/gwmap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "dubbo.apache.org/dubbo-go/v3/oteldubbo"
)

func init() {
	extension.SetFilter(keys.OtelClientTraceFilterKey, func() filter.Filter {
		return &otelClientFilter{}
	})
}

type otelClientFilter struct {
}

func (f *otelClientFilter) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	return result
}

func (f *otelClientFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	tracer := otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(dubbogoTrace.SemVersion()),
	)

	var span trace.Span
	ctx, span = tracer.Start(
		ctx,
		invocation.ActualMethodName(),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			dubbogoTrace.RPCSystemDubbo,
			semconv.RPCServiceKey.String(invoker.GetURL().ServiceKey()),
			semconv.RPCMethodKey.String(invocation.MethodName()),
		),
	)
	defer span.End()

	attachments := invocation.Attachments()
	if attachments == nil {
		attachments = map[string]interface{}{}
	}
	otel.GetTextMapPropagator().Inject(ctx, &gwmap.StrStrMap{InnerMap: attachments})
	for k, v := range attachments {
		invocation.SetAttachment(k, v)
	}
	result := invoker.Invoke(ctx, invocation)

	if result.Error() != nil {
		span.SetStatus(codes.Error, result.Error().Error())
	} else {
		span.SetStatus(codes.Ok, codes.Ok.String())
	}
	return result
}
