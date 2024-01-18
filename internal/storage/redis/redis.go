package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(address string, port string, password string, database int) *RedisStorage {
	return &RedisStorage{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", address, port),
			Password: password,
			DB:       database,
		}),
	}
}

func (s *RedisStorage) GetCounter(key string) (int, error) {
	countStr, err := s.client.Get(context.Background(), "ratelimits:req_qnt:"+key).Result()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *RedisStorage) IncrementCounter(key string, ttl int) error {
	newKey, err := s.client.Exists(context.Background(), "ratelimits:req_qnt:").Result()

	if err != nil {
		return err
	}

	err = s.client.Incr(context.Background(), "ratelimits:req_qnt:"+key).Err()

	if err != nil {
		return err
	}

	if newKey == 0 {
		err := s.client.Expire(context.Background(), key, time.Duration(ttl)*time.Second).Err()

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RedisStorage) RegisterBlock(key string, cooldown int) error {
	err := s.client.Set(context.Background(), "ratelimits:blocked:"+key, cooldown, time.Second*time.Duration(cooldown)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) IsBlocked(key string) (bool, error) {
	exists, err := s.client.Exists(context.Background(), "ratelimits:blocked:"+key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
