package mocks

import (
	"context"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
	"github.com/stretchr/testify/mock"
)

var _ storage.Storage = (*StorageMock)(nil)

type StorageMock struct {
	mock.Mock
}

func NewStorageMock() *StorageMock {
	return &StorageMock{}
}

func (sm *StorageMock) Increment(ctx context.Context, key string, ttl int) (int, error) {
	args := sm.Called(ctx, key, ttl)
	return args.Int(0), args.Error(1)
}

func (sm *StorageMock) Get(ctx context.Context, key string) (interface{}, error) {
	args := sm.Called(ctx, key)
	return args.Get(0), args.Error(1)
}

func (sm *StorageMock) Set(ctx context.Context, key string, ttl int) error {
	args := sm.Called(ctx, key, ttl)
	return args.Error(0)
}

func (sm *StorageMock) Exists(ctx context.Context, key string) (bool, error) {
	args := sm.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}
