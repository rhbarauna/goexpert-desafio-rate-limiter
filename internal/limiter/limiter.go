package limiter

import (
	"errors"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/configs"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
)

type Limiter struct {
	storage        storage.Storage
	cooldown       int
	ipDefaultLimit int
	ipTTL          int
	tokens         map[string]configs.TokenConfig
}

var ErrLimitedAccess = errors.New("access blocked")

func NewLimiter(storage storage.Storage, cooldown int, ipDefaultLimit int, ipTTL int, tokens map[string]configs.TokenConfig) Limiter {
	limiter := Limiter{
		storage:        storage,
		cooldown:       cooldown,
		ipDefaultLimit: ipDefaultLimit,
		ipTTL:          ipTTL,
		tokens:         tokens,
	}

	return limiter
}

func (l *Limiter) Limit(ip string, token string) error {
	limit := l.ipDefaultLimit
	ttl := l.ipTTL
	term := ip

	tokenConfig, tokenExists := l.tokens[token]

	if tokenExists {
		term = tokenConfig.Name
		ttl = tokenConfig.Ttl
		limit = tokenConfig.MaxRequests
	}

	return l.goNext(term, limit, ttl)
}

func (l *Limiter) goNext(term string, limit int, reqTTL int) error {
	counter, err := l.storage.GetCounter(term)

	if err != nil {
		return err
	}

	if counter == limit {
		l.storage.RegisterBlock(term, l.cooldown)
		return ErrLimitedAccess
	}

	l.storage.IncrementCounter(term, reqTTL)
	return nil
}
