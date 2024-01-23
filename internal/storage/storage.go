package storage

import (
	"context"
)

type Storage interface {
	Increment(ctx context.Context, key string, ttl int) (int, error)
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, ttl int) error
	Exists(ctx context.Context, key string) (bool, error)
}
