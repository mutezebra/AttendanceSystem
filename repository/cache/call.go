package cache

import (
	"context"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type CallCache struct {
	c *redis.Client
}

func NewCallCache() *CallCache {
	return &CallCache{
		c: RedisClient,
	}
}

func convert(uids []int64) []interface{} {
	interS := make([]interface{}, len(uids))
	for i := range uids {
		interS[i] = uids[i]
	}
	return interS
}

func (cache *CallCache) SetNewEvent(ctx context.Context, classID int64, uids []int64) error {
	_ = cache.DelCallEventSet(ctx, classID)
	if err := cache.c.SAdd(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID), convert(uids)...).Err(); err != nil {
		return errors.Wrap(err, "failed to set undo event")
	}
	return nil
}

func (cache *CallCache) AddDoneUser(ctx context.Context, classID int64, uids []int64) error {
	if err := cache.c.SAdd(ctx, fmt.Sprintf("%s:%d", consts.CallEventDoKey, classID), convert(uids)...).Err(); err != nil {
		return errors.Wrap(err, "failed to add done user")
	}
	return nil
}

func (cache *CallCache) WhetherUserExist(ctx context.Context, classID, uid int64) (exist bool, err error) {
	if exist, err = cache.c.SIsMember(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID), uid).Result(); err != nil {
		return false, errors.Wrap(err, "failed to check whether user exist")
	}
	return exist, nil
}

func (cache *CallCache) WhetherUserHaveDone(ctx context.Context, classID, uid int64) (exist bool, err error) {
	if exist, err = cache.c.SIsMember(ctx, fmt.Sprintf("%s:%d", consts.CallEventDoKey, classID), uid).Result(); err != nil {
		return false, errors.Wrap(err, "failed to check whether user have done")
	}
	return exist, nil
}

// GetUidDiff returns the difference between the undo set and the do set.也就是获得还未做的人
func (cache *CallCache) GetUidDiff(ctx context.Context, classID int64) ([]int64, error) {
	var s []int64
	if err := cache.c.SDiff(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID), fmt.Sprintf("%s:%d", consts.CallEventDoKey, classID)).ScanSlice(&s); err != nil {
		return nil, errors.Wrap(err, "failed to get uid diff")
	}
	return s, nil
}

// GetUidInter returns the intersection between the undo set and the do set.也就是获得已经做了的人
func (cache *CallCache) GetUidInter(ctx context.Context, classID int64) ([]int64, error) {
	var s []int64
	if err := cache.c.SInter(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID), fmt.Sprintf("%s:%d", consts.CallEventDoKey, classID)).ScanSlice(&s); err != nil {
		return nil, errors.Wrap(err, "failed to get uid inter")
	}
	return s, nil
}

func (cache *CallCache) DelCallEventSet(ctx context.Context, classID int64) error {
	if err := cache.c.Del(ctx, fmt.Sprintf("%s:%d", consts.CallEventUndoKey, classID), fmt.Sprintf("%s:%d", consts.CallEventDoKey, classID)).Err(); err != nil {
		return errors.Wrap(err, "failed to del call event set")
	}
	return nil
}

func (cache *CallCache) SaveSvcCallEvent(ctx context.Context, classID int64, data []byte, expire time.Duration) error {
	if err := cache.c.Set(ctx, fmt.Sprintf("%s:%d", consts.SvcCallEventKey, classID), data, expire).Err(); err != nil {
		return errors.Wrap(err, "failed to save svc call event")
	}
	return nil
}

func (cache *CallCache) ReadSvcCallEvent(ctx context.Context) ([][]byte, error) {
	var cursor uint64
	var res [][]byte

	for {
		keys, nextCursor, err := cache.c.Scan(ctx, cursor, fmt.Sprintf("%s:*", consts.SvcCallEventKey), 100).Result()
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan svc call event keys")
		}

		for _, key := range keys {
			data, err := cache.c.Get(ctx, key).Bytes()
			if err != nil {
				return nil, errors.Wrap(err, "failed to read svc call event")
			}
			res = append(res, data)
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return res, nil
}
