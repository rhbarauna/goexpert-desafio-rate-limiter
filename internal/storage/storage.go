package storage

import "github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/database"

type Storage interface {
	SetTokenConfig(token string, rateLimitInfo database.RateLimitInfo) error
	GetTokenConfig(token string) (database.RateLimitInfo, error)
	GetCounter(key string) (int, error)
	IncrementCounter(key string, ttl int) error
	RegisterBlock(key string, cooldown int) error
	IsBlocked(key string) (bool, error)
}
