package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"rag-new/internal/base/conf"
	"strconv"
)

func NewRedis(c *conf.Config) *redis.Client {
	//var pass string
	var r = redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	status := r.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}

	return r
}
