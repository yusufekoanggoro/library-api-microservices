package redis

import (
	"book-service/config"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Ping(ctx context.Context) error
	GetClient() *redis.Client
}

type redisService struct {
	client *redis.Client
}

func NewRedisService(cfg config.ConfigProvider) RedisClient {

	db, err := strconv.Atoi(cfg.GetRedisDB())
	if err != nil {
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.GetRedisHost(), cfg.GetRedisPort()),
		Password: cfg.GetRedisPassowrd(),
		DB:       db,
	})

	return &redisService{client: client}
}

func (r *redisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisService) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisService) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisService) GetClient() *redis.Client {
	return r.client
}
