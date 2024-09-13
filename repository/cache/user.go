package cache

import (
	"context"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	c *redis.Client
}

func NewUserCache() *UserCache {
	return &UserCache{
		c: RedisClient,
	}
}

func (cache *UserCache) WhetherVerifyCodeExist(ctx context.Context, phoneNumber string) (bool, string, error) {
	v, err := cache.c.Get(ctx, phoneNumber).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, "", errors.Wrap(err, fmt.Sprintf("failed when find %s`s verify code", phoneNumber))
	}

	return !errors.Is(err, redis.Nil), v, nil
}

func (cache *UserCache) PutVerifyCode(ctx context.Context, phoneNumber, verifyCode string) error {
	err := cache.c.Set(ctx, phoneNumber, verifyCode, consts.VerifyCodeExpireTime).Err()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed when set %s`s verifyCode", phoneNumber))
	}
	return nil
}
