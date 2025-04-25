package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Locker interface {
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string)
}

type RedisLocker struct {
	client *redis.Client
}

func NewRedisLocker(client *redis.Client) *RedisLocker {
	return &RedisLocker{client: client}
}

func (rl *RedisLocker) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	return rl.client.SetNX(ctx, lockKey, "locked", ttl).Result()
}

func (rl *RedisLocker) ReleaseLock(ctx context.Context, key string) {
	lockKey := fmt.Sprintf("lock:%s", key)
	rl.client.Del(ctx, lockKey)
}
