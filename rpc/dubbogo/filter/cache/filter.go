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

package cache

import (
	"context"
	"reflect"
	"strings"

	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/wesleywu/gowing/rpc/dubbogo/keys"
	"github.com/wesleywu/gowing/util/gwcache"
	"github.com/wesleywu/gowing/util/gwreflect"
	"google.golang.org/protobuf/proto"
)

const (
	cachedResult     = "cached-result"
	downgradedResult = "downgraded-result"
	cacheTagName     = "cache"
)

func init() {
	extension.SetFilter(keys.CacheFilterKey, newFilter)
}

func newFilter() filter.Filter {
	if !gwcache.CacheEnabled {
		return &cacheFilter{cacheEnabled: false}
	}
	ctx := gctx.New()
	cache, err := gwcache.GetCacheProvider()
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return &cacheFilter{cacheEnabled: false}
	}
	return &cacheFilter{
		cacheEnabled: true,
		cache:        cache,
	}
}

type cacheFilter struct {
	cacheEnabled bool
	cache        gwcache.CacheProvider
}

func (f *cacheFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	ctx, span := gtrace.NewSpan(ctx, "cacheFilter.Invoke")
	defer span.End()
	for f.cacheEnabled {
		// get providerUrl. The origin url may be is registry URL.
		invokerUrl := invoker.GetURL()
		providerUrl := getProviderURL(invokerUrl)
		service := invokerUrl.ServiceKey()

		params := invocation.Arguments()
		methodName := invocation.ActualMethodName()
		if len(params) != 1 {
			g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter. %s.%s(%s params)", service, methodName, len(params))
			break
		}
		req, ok := params[0].(proto.Message)
		if req == nil {
			g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter of type proto.Message. %s.%s(nil)", service, methodName)
			break
		}
		if !ok {
			g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter of type proto.Message. %s.%s(%s)", service, methodName, reflect.TypeOf(req).Kind())
			break
		}

		resultType, err := getMethodResultType(ctx, providerUrl, methodName)
		if err != nil {
			g.Log().Warningf(ctx, "cache filter failed to determine the return type of RPC method %s.%s(%s)", service, methodName, reflect.TypeOf(req).Kind())
			break
		}
		if resultType.Kind() != reflect.Ptr {
			g.Log().Warningf(ctx, "The return type of RPC method %s.%s should be a pointer, but a %s", service, methodName, resultType.Kind())
			break
		}
		result, err := f.getCachedResult(ctx, service, methodName, req, resultType)
		if err != nil { // error occurred when getting cached result
			if err == gwcache.ErrLockTimeout { // 获取锁超时，返回降级的结果
				invocation.SetAttachment(downgradedResult, true)
				return &protocol.RPCResult{
					Attrs: invocation.Attachments(),
					Err:   nil,
					Rest:  result,
				}
			} else if err == gwcache.ErrNotFound { // cache未找到
				break
			}
			// 其他底层错误
			g.Log().Errorf(ctx, "%+v", err)
			break
		}
		if result != nil { // cached result exists
			if gwcache.DebugEnabled {
				g.Log().Debugf(ctx, "rpc method\n\t%s.%s('%s') returned %s cached result\n\t%s", service, methodName, gjson.MustEncode(req), gwcache.CacheProviderName, gjson.MustEncodeString(result))
			}
			invocation.SetAttachment(cachedResult, true)
			return &protocol.RPCResult{
				Attrs: invocation.Attachments(),
				Err:   nil,
				Rest:  result,
			}
		}
		break
	}
	// no cached result, proceed to RPC call
	return invoker.Invoke(ctx, invocation)
}

func (f *cacheFilter) OnResponse(ctx context.Context, res protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	ctx, span := gtrace.NewSpan(ctx, "cacheFilter.OnResponse")
	defer span.End()
	if f.cacheEnabled {
		if gconv.Bool(res.Attachment(cachedResult, false)) { // is a cached res
			return res
		}
		if res != nil && res.Error() == nil { // normal res without error, may need to save it back to cache
			// get providerUrl. The origin url may be is registry URL.
			invokerUrl := invoker.GetURL()
			service := invokerUrl.ServiceKey()

			params := invocation.Arguments()
			methodName := invocation.ActualMethodName()

			if len(params) != 1 {
				g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter. %s.%s(%s params)", service, methodName, len(params))
				return res
			}
			req, ok := params[0].(proto.Message)
			if req == nil {
				g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter of type proto.Message. %s.%s(nil)", service, methodName)
				return res
			}
			if !ok {
				g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter of type proto.Message. %s.%s(%s)", service, methodName, reflect.TypeOf(req).Kind())
				return res
			}

			err := f.saveCachedResult(ctx, service, methodName, req, res.Result())
			if err != nil {
				g.Log().Error(ctx, "Error saving response back to cache: ", err)
			} else {
				if gwcache.DebugEnabled {
					g.Log().Debugf(ctx, "saved res of rpc method\n\t%s.%s('%s') to %s cache\n\t%s", service, methodName, gjson.MustEncode(params[0]), gwcache.CacheProviderName, gjson.MustEncodeString(res.Result()))
				}
			}
		}
	}
	return res
}

func getProviderURL(url *common.URL) *common.URL {
	if url.SubURL == nil {
		return url
	}
	return url.SubURL
}

func getMethodResultType(_ context.Context, url *common.URL, methodName string) (reflect.Type, error) {
	urlProtocol := url.Protocol
	path := strings.TrimPrefix(url.Path, "/")

	// get service
	svc := common.ServiceMap.GetServiceByServiceKey(urlProtocol, url.ServiceKey())
	if svc == nil {
		return nil, gerror.Newf("cannot find service [%s] in %s", path, urlProtocol)
	}

	// get method
	method := svc.Method()[methodName]
	if method == nil {
		return nil, gerror.Newf("cannot find method [%s] of service [%s] in %s", methodName, path, urlProtocol)
	}
	return method.ReplyType(), nil
}

func (f *cacheFilter) getCachedResult(ctx context.Context, service, methodName string, req proto.Message, resultType reflect.Type) (interface{}, error) {
	ctx, span := gtrace.NewSpan(ctx, "getCachedResult")
	defer span.End()
	metaField, err := gwreflect.GetMetaField(ctx, req)
	if err != nil {
		return nil, err
	}
	if metaField == nil {
		return nil, nil
	}

	cacheSetting := gwcache.ParseCacheTag(ctx, metaField.Tag.Get(cacheTagName))
	if !cacheSetting.Enabled {
		return nil, nil
	}

	err = gwreflect.MergeDefaultStructValue(ctx, req)
	if err != nil {
		return nil, err
	}
	if f.cache == nil || !f.cache.Initialized() {
		g.Log().Warning(ctx, "cacheFilter invoke error: cache not initialized")
		return nil, nil
	}
	// instantiate result as a pointer of result type
	result, ok := reflect.New(resultType.Elem()).Interface().(proto.Message)
	if !ok {
		return nil, gerror.Newf("The return type of RPC method %s.%s should be a protobuf message", service, methodName)
	}
	cacheKey := gwcache.GetCacheKey(ctx, service, methodName, req)
	if cacheKey == nil {
		return nil, err
	}
	err = f.cache.RetrieveCacheTo(ctx, cacheKey, result)
	// 返回缓存的结果
	return result, err
}

func (f *cacheFilter) saveCachedResult(ctx context.Context, service, methodName string, req proto.Message, result interface{}) error {
	ctx, span := gtrace.NewSpan(ctx, "saveCachedResult")
	defer span.End()
	if result == nil || f.cache == nil || !f.cache.Initialized() {
		return nil
	}
	metaField, err := gwreflect.GetMetaField(ctx, req)
	if err != nil {
		return err
	}
	if metaField == nil {
		return nil
	}
	cacheSetting := gwcache.ParseCacheTag(ctx, metaField.Tag.Get(cacheTagName))
	if !cacheSetting.Enabled {
		return nil
	}
	ttlSeconds := cacheSetting.Ttl

	cacheKey := gwcache.GetCacheKey(ctx, service, methodName, req)
	if cacheKey != nil {
		return f.cache.SaveCache(ctx, service, cacheKey, result, ttlSeconds)
	}
	return nil
}
