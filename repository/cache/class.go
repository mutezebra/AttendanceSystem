package cache

import (
	"context"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/pkg/errors"
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

func (cache *ClassCache) GetUID(ctx context.Context) int64 {
	return cache.c.Incr(ctx, consts.UIDCacheKey).Val()
}

func (cache *ClassCache) WhetherEventExist(ctx context.Context, classID int64) (exist bool, err error) {
	var existInt int64
	if existInt, err = cache.c.Exists(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID)).Result(); err != nil {
		return false, errors.Wrap(err, "failed to check whether event exist")
	}
	return existInt > 0, nil
}

// GetUidDiff returns the difference between the undo set and the do set.也就是获得还未做的人
func (cache *ClassCache) GetUidDiff(ctx context.Context, classID int64) ([]int64, error) {
	var s []int64
	if err := cache.c.SDiff(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID), fmt.Sprintf("%s:%d", consts.CallEventDoKey, classID)).ScanSlice(&s); err != nil {
		return nil, errors.Wrap(err, "failed to get uid diff")
	}
	return s, nil
}

// GetUidInter returns the intersection between the undo set and the do set.也就是获得已经做了的人
func (cache *ClassCache) GetUidInter(ctx context.Context, classID int64) ([]int64, error) {
	var s []int64
	if err := cache.c.SInter(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID), fmt.Sprintf("%s:%d", consts.CallEventDoKey, classID)).ScanSlice(&s); err != nil {
		return nil, errors.Wrap(err, "failed to get uid inter")
	}
	return s, nil
}
