//go:build wireinject
// +build wireinject

package main

import (
	"path/filepath"
	"runtime"

	"github.com/google/wire"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/configs"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/limiter"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/middleware"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage/redis"
)

func provideConfig() *configs.Config {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("Erro ao obter informações do arquivo.")
	}
	goDir := filepath.Dir(currentFile)

	config, err := configs.LoadConfig(goDir)
	if err != nil {
		panic(err)
	}
	return config
}

func provideStorage(config *configs.Config) storage.Storage {
	return redis.NewRedisStorage(config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
}

func provideLimiter(storage storage.Storage, config *configs.Config) limiter.LimiterInterface {
	return limiter.NewLimiter(storage, config.Cooldown, config.MaxRequests, config.Ttl, config.Tokens)
}

func provideRateLimiter(limiter limiter.LimiterInterface) middleware.RateLimiter {
	return middleware.NewRateLimiter(limiter)
}

func GetWebServerPort() string {
	return provideConfig().WebServerPort
}

func NewRateLimiter() (middleware.RateLimiter, error) {
	wire.Build(
		provideConfig,
		provideStorage,
		provideLimiter,
		provideRateLimiter,
	)
	return middleware.RateLimiter{}, nil
}
