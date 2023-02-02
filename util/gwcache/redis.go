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
	"fmt"
	"github.com/bsm/redislock"
	"github.com/cespare/xxhash/v2"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
	"time"
)

const (
	defaultPoolSize        = 40
	defaultPoolMinIdle     = 5
	defaultPoolMaxIdle     = 10
	defaultPoolMaxActive   = 100
	defaultPoolIdleTimeout = 10 * time.Second
	defaultPoolWaitTimeout = 10 * time.Second
	defaultPoolMaxLifeTime = 30 * time.Second
	defaultMaxRetries      = -1
)

var (
	syncLockTimeout     time.Duration
	redisNilResultError = redis.Nil
)

type redisCacheProvider struct {
	cacheInitialized bool
	client           redis.UniversalClient
	locker           *redislock.Client
}

func NewRedisCacheProvider(ctx context.Context) (CacheProvider, error) {
	var (
		cacheInitialized bool
		redisClient      redis.UniversalClient
		redisLocker      *redislock.Client
	)
	config, ok := gredis.GetConfig()
	if !ok {
		config, _ = gredis.ConfigFromMap(g.Map{
			"Address": "127.0.0.1:6379",
		})
	}
	fillWithDefaultConfiguration(config)
	opts := &redis.UniversalOptions{
		Addrs:           gstr.SplitAndTrim(config.Address, ","),
		Password:        config.Pass,
		DB:              config.Db,
		MaxRetries:      defaultMaxRetries,
		PoolSize:        defaultPoolSize,
		MinIdleConns:    config.MinIdle,
		MaxIdleConns:    config.MaxIdle,
		ConnMaxLifetime: config.MaxConnLifetime,
		ConnMaxIdleTime: config.IdleTimeout,
		PoolTimeout:     config.WaitTimeout,
		DialTimeout:     config.DialTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
		MasterName:      config.MasterName,
		TLSConfig:       config.TLSConfig,
	}

	if opts.MasterName != "" {
		redisSentinel := opts.Failover()
		redisSentinel.ReplicaOnly = config.SlaveOnly
		redisClient = redis.NewFailoverClient(redisSentinel)
	} else if len(opts.Addrs) > 1 {
		redisClient = redis.NewClusterClient(opts.Cluster())
	} else {
		redisClient = redis.NewClient(opts.Simple())
	}
	if TracingEnabled {
		if err := redisotel.InstrumentTracing(redisClient); err != nil {
			g.Log().Errorf(ctx, "failed to trace redis via otel: %+v", err)
		}
	}

	info := redisClient.Info(ctx)
	if info.Err() != nil {
		return nil, info.Err()
	}
	if redisClient != nil {
		redisLocker = redislock.New(redisClient)
	}

	lockTimeoutMillisVar, err := g.Cfg().Get(ctx, "cache.syncTimeoutMillis", 100)
	lockTimeoutMillis := 100
	if err == nil {
		lockTimeoutMillis = lockTimeoutMillisVar.Int()
	}
	syncLockTimeout = time.Duration(lockTimeoutMillis) * time.Millisecond
	g.Log().Info(ctx, "redis cache provider initialized")
	if DebugEnabled {
		g.Log().Debug(ctx, info.String())
	}
	cacheInitialized = true

	return &redisCacheProvider{
		cacheInitialized: cacheInitialized,
		client:           redisClient,
		locker:           redisLocker,
	}, nil
}

func (s *redisCacheProvider) Initialized() bool {
	return s.cacheInitialized
}

// GetCacheKey 生成cacheKey
// serviceName service名称，不同service名称不要相同，否则会造成 cacheKey 冲突，可以将 serviceName 当做某些缓存实现的 namespace 看待
// funcName method名称，不同funcName名称不要相同，否则会造成 cacheKey 冲突
// funcParams 所有的method参数
func GetCacheKey(ctx context.Context, serviceName string, funcName string, funcParams proto.Message) *string {
	ctx, span := gtrace.NewSpan(ctx, "GetCacheKey_"+serviceName+"_"+funcName)
	defer span.End()
	cacheKey := ServiceCachePrefix + serviceName + "_" + funcName
	if funcParams != nil {
		paramBytes, err := proto.Marshal(funcParams)
		if err != nil {
			paramBytes = gjson.MustEncode(funcParams)
		}
		hash := xxhash.Sum64(paramBytes)
		cacheKey = fmt.Sprintf("%s:%x", cacheKey, hash)
	}
	return &cacheKey
}

// RetrieveCacheTo 根据 cacheKey 获取缓存对象，并通过 json 解码到 value 中
// value 应该是原始对象的指针，必须在外部先初始化该对象
func (s *redisCacheProvider) RetrieveCacheTo(ctx context.Context, cacheKey *string, value proto.Message) error {
	ctx, span := gtrace.NewSpan(ctx, "RetrieveCacheTo"+*cacheKey)
	defer span.End()
	if cacheKey == nil {
		return ErrEmptyCacheKey
	}
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	var (
		lock *redislock.Lock
		err  error
	)
	_, spanLock := gtrace.NewSpan(ctx, "LockCache"+*cacheKey)
	// 对每次取 cache  给最长 3秒钟处理时间，在此期间到达的同样取 cache 请求会等待当前处理结束
	lock, err = s.locker.Obtain(ctx, ServiceCacheLockerPrefix+*cacheKey, syncLockTimeout, &redislock.Options{
		RetryStrategy: redislock.ExponentialBackoff(10*time.Millisecond, syncLockTimeout),
		Metadata:      "",
	})
	if err == redislock.ErrNotObtained {
		g.Log().Errorf(ctx, "Timeout obtaining lock for cache \"%s\" %s", *cacheKey, err.Error())
		spanLock.End()
		return ErrLockTimeout
	} else if err != nil {
		g.Log().Errorf(ctx, "Error obtaining lock for cache \"%s\" %s", *cacheKey, err.Error())
		spanLock.End()
		return err
	}
	defer func(lock *redislock.Lock, ctx context.Context) {
		_, spanLockRelease := gtrace.NewSpan(ctx, "LockReleaseCache"+*cacheKey)
		_ = lock.Release(ctx)
		spanLockRelease.End()
	}(lock, ctx)
	spanLock.End()

	cachedValue, err := s.client.Get(ctx, *cacheKey).Bytes()
	if err == redisNilResultError {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	_, unmarshal := gtrace.NewSpan(ctx, "Unmarshal"+*cacheKey)
	err = deserialize(cachedValue, value)
	if err != nil {
		g.Log().Warningf(ctx, "error decoding cache value \"%s\" for key \"%s\" %s", cachedValue, *cacheKey, err.Error())
		unmarshal.End()
		return err
	}
	unmarshal.End()
	return nil
}

// SaveCache 添加缓存
// serviceName service名称，需要的原因是要记录当前 service 下所有已经保存的缓存 key 的集合
// cacheKey 缓存key
// value 要放入缓存的value，保存前会对其进行 json 编码
func (s *redisCacheProvider) SaveCache(ctx context.Context, serviceName string, cacheKey *string, value any, ttlSeconds uint32) error {
	ctx, span := gtrace.NewSpan(ctx, "SaveCache"+*cacheKey)
	defer span.End()
	if cacheKey == nil {
		return ErrEmptyCacheKey
	}
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	valueBytes, err := serialize(value.(proto.Message))
	if err != nil {
		g.Log().Warningf(ctx, "error encoding cache value for key \"%s\" %s", *cacheKey, err.Error())
		return err
	}
	cacheKeysetName := ServiceCacheKeySetPrefix + serviceName

	pipe := s.client.Pipeline()
	pipe.Set(ctx, *cacheKey, valueBytes, time.Duration(ttlSeconds)*time.Second)
	pipe.SAdd(ctx, cacheKeysetName, *cacheKey)
	_, err = pipe.Exec(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "error save cache key \"%s\" to keyset \"%s\" %s", *cacheKey, cacheKeysetName, err.Error())
		return err
	}
	return nil
}

// RemoveCache 根据 cacheKey 删除单个缓存
// cacheKey 缓存key
func (s *redisCacheProvider) RemoveCache(ctx context.Context, cacheKey *string) error {
	if cacheKey == nil {
		return ErrEmptyCacheKey
	}
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	err := s.client.Del(ctx, *cacheKey).Err()
	if err != nil {
		g.Log().Warningf(ctx, "error delete cache key \"%s\" %s", cacheKey, err.Error())
		return err
	}
	return nil
}

// ClearCache 清除 service 下所有已经保存的缓存
func (s *redisCacheProvider) ClearCache(ctx context.Context, serviceName string) error {
	if !s.Initialized() {
		return ErrCacheNotInitialized
	}
	cacheKeysetName := ServiceCacheKeySetPrefix + serviceName
	keys, err := s.client.SMembers(ctx, cacheKeysetName).Result()
	if err != nil {
		g.Log().Warningf(ctx, "error load cache keyset \"%s\" %s", cacheKeysetName, err.Error())
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	err = s.client.Del(ctx, append(keys, cacheKeysetName)...).Err()
	if err != nil {
		g.Log().Warningf(ctx, "error clear cache keys \"%s\" %s", gstr.Join(keys, ","), err.Error())
		return err
	}
	return nil
}

func fillWithDefaultConfiguration(config *gredis.Config) {
	if config.MinIdle == 0 {
		config.MinIdle = defaultPoolMinIdle
	}
	// The MaxIdle is the most important attribute of the connection pool.
	// Only if this attribute is set, the created connections from client
	// can not exceed the limit of the server.
	if config.MaxIdle == 0 {
		config.MaxIdle = defaultPoolMaxIdle
	}
	// This value SHOULD NOT exceed the connection limit of redis server.
	if config.MaxActive == 0 {
		config.MaxActive = defaultPoolMaxActive
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = defaultPoolIdleTimeout
	}
	if config.WaitTimeout == 0 {
		config.WaitTimeout = defaultPoolWaitTimeout
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = defaultPoolMaxLifeTime
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = -1
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = -1
	}
}
