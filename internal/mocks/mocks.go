package mocks

import (
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
	"github.com/stretchr/testify/mock"
)

var _ storage.Storage = (*StorageMock)(nil)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) GetCounter(key string) (int, error) {
	args := sm.Called(key)
	return args.Int(0), args.Error(1)
}

func (sm *StorageMock) IncrementCounter(key string, ttl int) (int64, error) {
	args := sm.Called(key, ttl)
	return int64(args.Int(0)), args.Error(1)
}

func (sm *StorageMock) RegisterBlock(key string, cooldown int) error {
	args := sm.Called(key, cooldown)
	return args.Error(0)
}

func (sm *StorageMock) IsBlocked(key string) (bool, error) {
	args := sm.Called(key)
	return args.Bool(0), args.Error(1)
}
