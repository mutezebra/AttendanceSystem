package cache

import (
	"context"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/log"

	redis "github.com/redis/go-redis/v9"

	"github.com/mutezebra/ClassroomRandomRollCallSystem/config"
)

var RedisClient *redis.Client

func InitCache() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Host + ":" + config.Conf.Redis.Port,
		Password: config.Conf.Redis.Password,
		Network:  config.Conf.Redis.Network,
		DB:       config.Conf.Redis.Database,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.LogrusObj.Panic(err)
	}
	RedisClient = client
}
