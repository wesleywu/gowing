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

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/wesleywu/gowing/rpc/inproc/types"
	"github.com/wesleywu/gowing/util/gwcache"
	"github.com/wesleywu/gowing/util/gwreflect"
	"google.golang.org/protobuf/proto"
)

const (
	cacheTagName = "cache"
)

type InterceptorCache[TReq proto.Message, TRes proto.Message] struct {
	cacheEnabled bool
	cache        gwcache.CacheProvider
}

func NewCacheInterceptor[TReq proto.Message, TRes proto.Message]() *InterceptorCache[TReq, TRes] {
	if !gwcache.CacheEnabled {
		return &InterceptorCache[TReq, TRes]{
			cacheEnabled: false,
		}
	}
	// todo support cache provider set by tag
	ctx := context.Background()
	cache, err := gwcache.GetCacheProvider()
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return &InterceptorCache[TReq, TRes]{
			cacheEnabled: false,
			cache:        cache,
		}
	}
	return &InterceptorCache[TReq, TRes]{
		cacheEnabled: true,
		cache:        cache,
	}
}

func (s *InterceptorCache[TReq, TRes]) InterceptMethod(next *types.MethodInvocation[TReq, TRes]) *types.MethodInvocation[TReq, TRes] {
	nextMethodFunc := func(ctx context.Context, req TReq) (result TRes, err error) {
		var (
			cacheSetting gwcache.Setting
			cacheKey     *string
			ok           bool
			resultType   reflect.Type
		)
		ctx, span := gtrace.NewSpan(ctx, "cacheInterceptor")
		defer span.End()
		metaField, err := gwreflect.GetMetaField(ctx, req)
		if err != nil {
			goto callNext
		}
		if metaField == nil {
			goto callNext
		}

		cacheSetting = gwcache.ParseCacheTag(ctx, metaField.Tag.Get(cacheTagName))
		if !cacheSetting.Enabled {
			goto callNext
		}

		err = gwreflect.MergeDefaultStructValue(ctx, req)
		if err != nil {
			goto callNext
		}
		if s.cache == nil || !s.cache.Initialized() {
			g.Log().Warning(ctx, "cacheFilter invoke error: cache not initialized")
			goto callNext
		}

		resultType = reflect.TypeOf(result)
		result, ok = reflect.New(resultType.Elem()).Interface().(TRes)
		if !ok {
			g.Log().Warning(ctx, "The return type of RPC method %s.%s should be a protobuf message", next.ServiceName, next.MethodName)
			goto callNext
		}
		cacheKey = gwcache.GetCacheKey(ctx, next.ServiceName, next.MethodName, req)
		if cacheKey == nil {
			goto callNext
		}
		err = s.cache.RetrieveCacheTo(ctx, cacheKey, result)
		// 返回缓存的结果
		if err != nil { // error occurred when getting cached result
			if err == gwcache.ErrLockTimeout { // 获取锁超时，返回降级的结果
				goto callNext
			} else if err == gwcache.ErrNotFound { // cache未找到
				goto callNext
			}
			// 其他底层错误
			g.Log().Errorf(ctx, "%+v", err)
			goto callNext
		}
		if result.ProtoReflect() != nil { // cached result exists
			if gwcache.DebugEnabled {
				g.Log().Debugf(ctx, "rpc method\n\t%s.%s('%s') returned %s cached result\n\t%s", next.ServiceName, next.MethodName, gjson.MustEncode(req), gwcache.CacheProviderName, gjson.MustEncodeString(result))
			}
			return result, nil
		}
	callNext:
		result, err = next.CallMethod(ctx, req)
		if cacheSetting.Enabled {
			if result.ProtoReflect() != nil && err == nil { // normal res without error, may need to save it back to cache
				if cacheKey == nil {
					return result, err
				}
				ttlSeconds := cacheSetting.Ttl
				err = s.cache.SaveCache(ctx, next.ServiceName, cacheKey, result, ttlSeconds)
				if err != nil {
					g.Log().Error(ctx, "Error saving response back to cache: ", err)
				} else {
					if gwcache.DebugEnabled {
						g.Log().Debugf(ctx, "saved res of rpc method\n\t%s.%s('%s') to %s cache\n\t%s", next.ServiceName, next.MethodName, gjson.MustEncode(req), gwcache.CacheProviderName, gjson.MustEncodeString(result))
					}
				}
			}
		}
		return result, err
	}
	return &types.MethodInvocation[TReq, TRes]{
		SvcMethod:   nextMethodFunc,
		ServiceName: next.ServiceName,
		MethodName:  next.MethodName,
	}
}
