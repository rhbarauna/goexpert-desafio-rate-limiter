package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
)

var _ storage.Storage = (*RedisStorage)(nil)

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

func (s *RedisStorage) IncrementCounter(key string, ttl int) (int64, error) {
	redisKey := fmt.Sprintf("ratelimits:req_qnt:%s", key)
	newKey, err := s.client.Exists(context.Background(), redisKey).Result()

	log.Printf("%s, is new? %v\n", redisKey, newKey)

	if err != nil {
		return 0, err
	}

	counter, err := s.client.Incr(context.Background(), redisKey).Result()

	if err != nil {
		return 0, err
	}

	// O método Exists retorna 1 se a chave existir, 0 se não existir
	if newKey == 0 {
		duration := time.Duration(ttl) * time.Second
		err := s.client.Expire(context.Background(), redisKey, duration).Err()

		if err != nil {
			log.Printf("erro ao setar o expire para a chave %s. %s\n", redisKey, err.Error())
			return 0, err
		}
		log.Printf("Setado o expire para a chave %s.\n", redisKey)
	}

	return counter, nil
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
