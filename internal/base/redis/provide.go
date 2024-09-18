package redis

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"rag-new/internal/base/conf"
	"strconv"
)

type Redis struct {
	Client *redis.Client
	Locker *redislock.Client
}

func NewRedis(c *conf.Config) *Redis {
	var client = redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	status := client.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}

	// Create a new lock client.
	locker := redislock.New(client)

	var r = &Redis{
		Client: client,
		Locker: locker,
	}

	return r
}
