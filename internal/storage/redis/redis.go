package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
)

var _ storage.Storage = (*redisStorage)(nil)

type redisStorage struct {
	client *redis.Client
}

func NewRedisStorage(address string, port string, password string, database int) *redisStorage {
	return &redisStorage{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", address, port),
			Password: password,
			DB:       database,
		}),
	}
}

func (s *redisStorage) Increment(ctx context.Context, key string, ttl int) (int, error) {
	pipe := s.client.Pipeline()

	pipe.Exists(ctx, key)
	pipe.Incr(ctx, key)

	counter, err := pipe.Exec(ctx) // Execute the pipeline

	if err != nil {
		return 0, err
	}

	if len(counter) > 0 && counter[0].(*redis.IntCmd).Val() == 0 { // Key didn't exist
		err = pipe.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	}

	if err != nil {
		return 0, err
	}

	_, err = pipe.Exec(ctx) // Ensure pipeline execution
	if err != nil {
		return 0, err
	}

	return int(counter[1].(*redis.IntCmd).Val()), nil
}

func (s *redisStorage) Get(ctx context.Context, key string) (interface{}, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *redisStorage) Set(ctx context.Context, key string, ttl int) error {
	var err error
	pipe := s.client.Pipeline()
	pipe.Set(ctx, key, true, 0)

	if ttl != 0 {
		err = pipe.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStorage) Exists(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}
