package memcache

import (
	"context"
	"first-proj/common"
	"time"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/redis/go-redis/v9"

	"github.com/go-redis/cache/v9"
)

// implements Cache interface
type redisCache struct {
	store *cache.Cache
}

func NewRedisCache(sc goservice.ServiceContext) *redisCache {
	rdClient := sc.MustGet(common.PluginRedis).(*redis.Client)

//	Use redis as the main caching store and a local in-memory caching store
	c := cache.New(&cache.Options{
		Redis:      rdClient,
		/* Stores 1000 items in memory for 1 minute.
			Returns a local cache instance that automatically removes
			least frequently used items when full. */
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	return &redisCache{store: c}
}

func (rdc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rdc.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   0, 
	})
}

func (rdc *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	return rdc.store.Get(ctx, key, value)
}

func (rdc *redisCache) Delete(ctx context.Context, key string) error {
	return rdc.store.Delete(ctx, key)
}