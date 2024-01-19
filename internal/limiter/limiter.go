package limiter

import (
	"errors"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/configs"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
)

type Limiter struct {
	storage     storage.Storage
	cooldown    int
	maxRequests int
	ttl         int
	tokens      map[string]configs.TokenConfig
}

var ErrLimitedAccess = errors.New("access blocked")

func NewLimiter(storage storage.Storage, cooldown int, maxRequests int, ttl int, tokens map[string]configs.TokenConfig) Limiter {
	limiter := Limiter{
		storage:     storage,
		cooldown:    cooldown,
		maxRequests: maxRequests,
		ttl:         ttl,
		tokens:      tokens,
	}

	return limiter
}

func (l *Limiter) Limit(ip string, token string) error {
	maxRequests := l.maxRequests
	term := ip
	cooldown := l.cooldown

	tokenConfig, tokenExists := l.tokens[token]

	if tokenExists {
		term = tokenConfig.Name
		maxRequests = tokenConfig.MaxRequests
		cooldown = tokenConfig.Cooldown
	}

	return l.goNext(term, maxRequests, cooldown)
}

func (l *Limiter) goNext(term string, maxRequests int, cooldown int) error {
	isBlocked, err := l.storage.IsBlocked(term)

	if err != nil {
		return err
	}

	if isBlocked {
		return ErrLimitedAccess
	}

	counter, err := l.storage.IncrementCounter(term, l.ttl)

	if err != nil {
		return err
	}

	if counter >= int64(maxRequests) {
		err = l.storage.RegisterBlock(term, cooldown)
	}

	if err != nil {
		return err
	}

	return nil
}
