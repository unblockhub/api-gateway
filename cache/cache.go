package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var (
	RedisClient *redis.Client
	Ctx = context.Background()
)

func Init()  {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"),
	})
}
