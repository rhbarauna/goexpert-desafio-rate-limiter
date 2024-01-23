package limiter

import (
	"context"
	"errors"
	"fmt"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/configs"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
)

var _ Limiter = (*limiter)(nil)

type Limiter interface {
	Limit(ip string, token string) error
}

type limiter struct {
	Storage     storage.Storage
	Cooldown    int
	MaxRequests int
	TTL         int
	Tokens      map[string]configs.TokenConfig
}

var ErrLimitedAccess = errors.New("access blocked")

func NewLimiter(storage storage.Storage, cooldown int, maxRequests int, ttl int, tokens map[string]configs.TokenConfig) *limiter {
	limiter := &limiter{
		Storage:     storage,
		Cooldown:    cooldown,
		MaxRequests: maxRequests,
		TTL:         ttl,
		Tokens:      tokens,
	}

	return limiter
}

func (l *limiter) Limit(ip string, token string) error {
	ctx := context.Background()
	term := ip
	maxRequests := l.MaxRequests
	cooldown := l.Cooldown

	if tokenConfig, ok := l.Tokens[token]; ok {
		maxRequests = tokenConfig.MaxRequests
		cooldown = tokenConfig.Cooldown
		term = token
	}

	reqQntKey := fmt.Sprintf("ratelimit:req_qnt:%s", term)
	blockedKey := fmt.Sprintf("ratelimit:blocked:%s", term)

	isBlocked, err := l.Storage.Exists(ctx, blockedKey)

	if err != nil {
		return err
	}

	if isBlocked {
		return ErrLimitedAccess
	}

	counter, err := l.Storage.Increment(ctx, reqQntKey, l.TTL)

	if err != nil {
		return err
	}

	if counter > maxRequests {
		err = l.Storage.Set(ctx, blockedKey, cooldown)
		if err != nil {
			return err
		}

		return ErrLimitedAccess
	}

	return nil
}
