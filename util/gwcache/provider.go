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

package gwcache

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"google.golang.org/protobuf/proto"
)

const (
	ProviderNameMemory       = "memory"
	ProviderNameRedis        = "redis"
	ServiceCachePrefix       = "_SC_"
	ServiceCacheKeySetPrefix = "_SC_SET_"
	ServiceCacheLockerPrefix = "_LOCK_"
)

var (
	ctx                    = gctx.New()
	CacheEnabled           = false
	DebugEnabled           = false
	TracingEnabled         = false
	CacheProviderName      = "unknown"
	cacheProviderMap       = gmap.StrAnyMap{}
	ErrCacheNotInitialized = errors.New("cache: not initialized")
	ErrEmptyCacheKey       = errors.New("cache: cache key is empty")
	ErrLockTimeout         = errors.New("cache: lock timeout")
	ErrNotFound            = errors.New("cache: not found")
	ErrEmptyCachedValue    = errors.New("cache: cached value is empty")
)

type CacheProvider interface {
	Initialized() bool
	RetrieveCacheTo(ctx context.Context, cacheKey *string, value proto.Message) error
	SaveCache(ctx context.Context, serviceName string, cacheKey *string, value any, ttlSeconds uint32) error
	RemoveCache(ctx context.Context, cacheKey *string) error
	ClearCache(ctx context.Context, serviceName string) error
}

func init() {
	cacheEnabledVar, err := g.Cfg().Get(ctx, "gwcache.enabled", false)
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return
	}
	cacheProviderVar, err := g.Cfg().Get(ctx, "gwcache.provider", "memory")
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return
	}
	debugEnabledVar, err := g.Cfg().Get(ctx, "gwcache.debug", false)
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return
	}
	tracingEnabledVar, err := g.Cfg().Get(ctx, "gwcache.tracing", false)
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return
	}
	CacheEnabled = cacheEnabledVar.Bool()
	DebugEnabled = debugEnabledVar.Bool()
	TracingEnabled = tracingEnabledVar.Bool()
	CacheProviderName = cacheProviderVar.String()
	if CacheEnabled {
		g.Log().Infof(ctx, "service cache enabled with provider '%s'", CacheProviderName)
	}
}

func createProviderIfNotExists(ctx context.Context, adapterType string, createProviderFunc func(context.Context) (CacheProvider, error)) (CacheProvider, error) {
	providerVar := cacheProviderMap.GetVarOrSetFuncLock(adapterType, func() interface{} {
		provider, err := createProviderFunc(ctx)
		if err != nil {
			return err
		}
		return provider
	})
	switch providerVar.Val().(type) {
	case error:
		return nil, providerVar.Val().(error)
	default:
		return providerVar.Val().(CacheProvider), nil
	}
}

func GetCacheProvider() (CacheProvider, error) {
	switch CacheProviderName {
	case ProviderNameMemory:
		return createProviderIfNotExists(ctx, ProviderNameMemory, NewMemoryCacheProvider)
	case ProviderNameRedis:
		return createProviderIfNotExists(ctx, ProviderNameRedis, NewRedisCacheProvider)
	default:
		return nil, gerror.Newf("Not supported cache provider: %s", CacheProviderName)
	}
}
