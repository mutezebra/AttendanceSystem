package cache

import (
	"context"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/redis/go-redis/v9"
)

type ClassCache struct {
	c *redis.Client
}

func NewClassCache() *ClassCache {
	return &ClassCache{
		c: RedisClient,
	}
}

func (cache *ClassCache) GetClassID(ctx context.Context) int64 {
	return cache.c.Incr(ctx, consts.ClassIDCacheKey).Val()
}
