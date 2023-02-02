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
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/WesleyWu/gowing/rpc/dubbogo/keys"
	"github.com/WesleyWu/gowing/util/gwcache"
	"github.com/WesleyWu/gowing/util/gwreflect"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
)

const (
	cachedResult     = "cached-result"
	downgradedResult = "downgraded-result"
	cacheTagName     = "cache"
	enabledDefault   = false
	adapterDefault   = "redis"
	ttlDefault       = uint32(600)
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

type Setting struct {
	Enabled bool
	Name    string // todo implements custom specified name
	Key     string // todo implements custom specified key
	Adapter string // todo implements memory adapter
	Ttl     uint32
}

func (f *cacheFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	ctx, span := gtrace.NewSpan(ctx, "cacheFilter.Invoke")
	defer span.End()
	if f.cacheEnabled {
		result, err := f.getCachedResult(ctx, invoker, invocation)
		if err != nil { // error occurred when getting cached result
			return &protocol.RPCResult{
				Attrs: invocation.Attachments(),
				Err:   err,
				Rest:  nil,
			}
		}
		if result != nil { // cached result exists
			if gwcache.DebugEnabled {
				params := invocation.Arguments()
				service := invoker.GetURL().ServiceKey()
				methodName := invocation.ActualMethodName()
				g.Log().Debugf(ctx, "rpc method\n\t%s.%s('%s') returned %s cached result\n\t%s", service, methodName, gjson.MustEncode(params[0]), gwcache.CacheProviderName, gjson.MustEncodeString(result))
			}
			invocation.SetAttachment(cachedResult, true)
			return &protocol.RPCResult{
				Attrs: invocation.Attachments(),
				Err:   nil,
				Rest:  result,
			}
		}
	}
	// no error or cached result, proceed to RPC call
	return invoker.Invoke(ctx, invocation)
}

func (f *cacheFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	ctx, span := gtrace.NewSpan(ctx, "cacheFilter.OnResponse")
	defer span.End()
	if f.cacheEnabled {
		if gconv.Bool(result.Attachment(cachedResult, false)) { // is a cached result
			return result
		}
		if result != nil && result.Error() == nil { // normal result without error, may need to save it back to cache
			err := f.saveCachedResult(ctx, invoker, invocation, result.Result())
			if err != nil {
				g.Log().Error(ctx, "Error saving response back to cache: ", err)
			} else {
				if gwcache.DebugEnabled {
					params := invocation.Arguments()
					service := invoker.GetURL().ServiceKey()
					methodName := invocation.ActualMethodName()
					g.Log().Debugf(ctx, "saved result of rpc method\n\t%s.%s('%s') to %s cache\n\t%s", service, methodName, gjson.MustEncode(params[0]), gwcache.CacheProviderName, gjson.MustEncodeString(result.Result()))
				}
			}
		}
	}
	return result
}

func parseCacheTag(ctx context.Context, cacheTag string) Setting {
	ctx, span := gtrace.NewSpan(ctx, "parseCacheTag")
	defer span.End()
	if cacheTag == "" {
		return Setting{
			Enabled: enabledDefault,
			Adapter: adapterDefault,
			Ttl:     ttlDefault,
		}
	}
	enabled := enabledDefault
	name := ""
	key := ""
	adapter := adapterDefault
	ttl := ttlDefault
	for _, s := range strings.Split(cacheTag, ",") {
		s = strings.Trim(s, " ")
		attrib := strings.Split(strings.Trim(s, " "), "=")
		if len(attrib) == 1 && attrib[0] == "true" {
			enabled = true
			continue
		}
		if len(attrib) == 2 {
			attribName := strings.Trim(attrib[0], " ")
			attribValue := strings.Trim(attrib[1], " ")
			switch attribName {
			case "name":
				name = attribValue
			case "key":
				key = attribValue
			case "adapter":
				adapter = attribValue
			case "ttl":
				ttl = gconv.Uint32(attribValue)
			}
		}
	}
	return Setting{
		Enabled: enabled,
		Name:    name,
		Key:     key,
		Adapter: adapter,
		Ttl:     ttl,
	}
}

func getProviderURL(url *common.URL) *common.URL {
	if url.SubURL == nil {
		return url
	}
	return url.SubURL
}

func getMethodResultType(_ context.Context, invoker protocol.Invoker, invocation protocol.Invocation) (reflect.Type, error) {
	// get providerUrl. The origin url may be is registry URL.
	url := getProviderURL(invoker.GetURL())

	methodName := invocation.MethodName()
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

func (f *cacheFilter) getCachedResult(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) (interface{}, error) {
	ctx, span := gtrace.NewSpan(ctx, "getCachedResult")
	defer span.End()
	params := invocation.Arguments()
	service := invoker.GetURL().ServiceKey()
	methodName := invocation.ActualMethodName()
	if len(params) != 1 {
		g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter. service name: %s, method Name: %s, param count: %s", service, methodName, len(params))
		return nil, nil
	}
	req, ok := params[0].(proto.Message)
	if req == nil {
		g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter of type proto.Message. service name: %s, method Name: %s, param is nil", service, methodName)
		return nil, nil
	}
	if !ok {
		g.Log().Warningf(ctx, "cacheFilter will only be applied on RPC method with only 1 non-nil parameter of type proto.Message. service name: %s, method Name: %s, param type: %s", service, methodName, reflect.TypeOf(req).Kind())
		return nil, nil
	}
	//g.Log().Debugf(ctx, "cacheFilter Invoke is called, service name: %s, method Name: %s, param value: %s", service, methodName, gjson.MustEncodeString(req))

	metaField, err := gwreflect.GetMetaField(ctx, req)
	if err != nil {
		return nil, err
	}
	if metaField == nil {
		return nil, nil
	}

	cacheSetting := parseCacheTag(ctx, metaField.Tag.Get(cacheTagName))
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
	methodResultType, err := getMethodResultType(ctx, invoker, invocation)
	if err != nil {
		return nil, err
	}
	if methodResultType.Kind() != reflect.Ptr {
		return nil, gerror.Newf("The return type of RPC method %s.%s should be a pointer", service, methodName)
	}
	// instantiate result as a pointer of result type
	result, ok := reflect.New(methodResultType.Elem()).Interface().(proto.Message)
	if !ok {
		return nil, gerror.Newf("The return type of RPC method %s.%s should be a protobuf message", service, methodName)
	}
	cacheKey := gwcache.GetCacheKey(ctx, service, methodName, req)
	if cacheKey == nil {
		return nil, err
	}
	err = f.cache.RetrieveCacheTo(ctx, cacheKey, result)
	if err != nil {
		if err == gwcache.ErrLockTimeout { // 获取锁超时，返回降级的结果
			invocation.SetAttachment(downgradedResult, true)
			return &protocol.RPCResult{
				Attrs: invocation.Attachments(),
				Err:   nil,
				Rest:  result,
			}, nil
		} else if err == gwcache.ErrNotFound { // cache 未找到，执行底层操作
			return nil, nil
		}
		// 其他底层错误
		g.Log().Errorf(ctx, "%+v", err)
		return nil, nil
	}
	// 返回缓存的结果
	return result, nil
}

func (f *cacheFilter) saveCachedResult(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation, result interface{}) error {
	ctx, span := gtrace.NewSpan(ctx, "saveCachedResult")
	defer span.End()
	if result == nil || f.cache == nil || !f.cache.Initialized() {
		return nil
	}
	params := invocation.Arguments()
	if len(params) != 1 {
		return nil
	}
	req, ok := params[0].(proto.Message)
	if req == nil {
		return nil
	}
	if !ok {
		return nil
	}
	metaField, err := gwreflect.GetMetaField(ctx, req)
	if err != nil {
		return err
	}
	if metaField == nil {
		return nil
	}
	cacheSetting := parseCacheTag(ctx, metaField.Tag.Get(cacheTagName))
	if !cacheSetting.Enabled {
		return nil
	}
	ttlSeconds := cacheSetting.Ttl

	service := invoker.GetURL().ServiceKey()
	methodName := invocation.ActualMethodName()
	cacheKey := gwcache.GetCacheKey(ctx, service, methodName, req)
	if cacheKey != nil {
		return f.cache.SaveCache(ctx, service, cacheKey, result, ttlSeconds)
	}
	return nil
}
