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
	"github.com/bsm/redislock"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/text/gstr"
	"google.golang.org/protobuf/proto"
	"time"
)

var memoryNilResultError = errors.New("nil result")

type memoryCacheProvider struct {
	cacheInitialized bool
	cache            *gcache.Cache
	locker           *redislock.Client
}

func NewMemoryCacheProvider(ctx context.Context) (CacheProvider, error) {
	cache := gcache.New()
	g.Log().Infof(ctx, "memory cache provider initialized")
	return &memoryCacheProvider{
		cacheInitialized: true,
		cache:            cache,
	}, nil
}

func (s *memoryCacheProvider) Initialized() bool {
	return s.cacheInitialized
}

func (s *memoryCacheProvider) RetrieveCacheTo(ctx context.Context, cacheKey *string, value proto.Message) error {
	if cacheKey == nil {
		return ErrEmptyCacheKey
	}
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	if TracingEnabled {
		var span *gtrace.Span
		ctx, span = gtrace.NewSpan(ctx, "RetrieveCacheTo "+*cacheKey)
		defer span.End()
	}
	cachedValue, err := s.cache.GetOrSetFuncLock(ctx, *cacheKey, func(ctx context.Context) (value interface{}, err error) {
		return nil, memoryNilResultError
	}, 0)
	if err == memoryNilResultError {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	_, unmarshal := gtrace.NewSpan(ctx, "Unmarshal"+*cacheKey)
	err = deserialize(cachedValue.Bytes(), value)
	if err != nil {
		g.Log().Warningf(ctx, "error decoding cache value \"%s\" for key \"%s\" %s", cachedValue, *cacheKey, err.Error())
		unmarshal.End()
		return err
	}
	unmarshal.End()
	return nil
}

func (s *memoryCacheProvider) SaveCache(ctx context.Context, serviceName string, cacheKey *string, value any, ttlSeconds uint32) error {
	if cacheKey == nil {
		return ErrEmptyCacheKey
	}
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	if TracingEnabled {
		var span *gtrace.Span
		ctx, span = gtrace.NewSpan(ctx, "SaveCache "+*cacheKey)
		defer span.End()
	}
	valueBytes, err := serialize(value.(proto.Message))
	if err != nil {
		g.Log().Warningf(ctx, "error encoding cache value for key \"%s\" %s", *cacheKey, err.Error())
		return err
	}
	cacheKeysetName := ServiceCacheKeySetPrefix + serviceName

	err = s.cache.Set(ctx, *cacheKey, valueBytes, time.Duration(ttlSeconds)*time.Second)
	if err != nil {
		return err
	}
	setVal, err := s.cache.Get(ctx, cacheKeysetName)
	if err != nil {
		return err
	}
	setValStrings := setVal.Strings()
	if gset.NewStrSetFrom(setValStrings, false).Contains(*cacheKey) {
		return nil
	}
	setValStrings = append(setValStrings, *cacheKey)
	return s.cache.Set(ctx, cacheKeysetName, setValStrings, 0)
}

func (s *memoryCacheProvider) RemoveCache(ctx context.Context, cacheKey *string) error {
	if cacheKey == nil {
		return ErrEmptyCacheKey
	}
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	if TracingEnabled {
		var span *gtrace.Span
		ctx, span = gtrace.NewSpan(ctx, "RemoveCache "+*cacheKey)
		defer span.End()
	}
	_, err := s.cache.Remove(ctx, *cacheKey)
	if err != nil {
		g.Log().Warningf(ctx, "error delete cache key \"%s\" %s", cacheKey, err.Error())
		return err
	}
	return nil
}

func (s *memoryCacheProvider) ClearCache(ctx context.Context, serviceName string) error {
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	if TracingEnabled {
		var span *gtrace.Span
		ctx, span = gtrace.NewSpan(ctx, "ClearCache "+serviceName)
		defer span.End()
	}
	cacheKeysetName := ServiceCacheKeySetPrefix + serviceName
	setVal, err := s.cache.Get(ctx, cacheKeysetName)
	if err != nil {
		return err
	}
	keys := setVal.Strings()
	if len(keys) == 0 {
		return nil
	}

	keysToRemote := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		keysToRemote[i] = key
	}
	keysToRemote[len(keys)] = cacheKeysetName

	_, err = s.cache.Remove(ctx, keysToRemote...)
	if err != nil {
		g.Log().Warningf(ctx, "error clear cache keys \"%s\" %s", gstr.Join(keys, ","), err.Error())
		return err
	}
	return nil
}
