package limiter

import (
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage"
)

type Limiter struct {
	Storage        storage.Storage
	IpCooldown     int
	IpDefaultLimit int
	IpTTL          int
}

func NewLimiter(storage storage.Storage, ipCooldown int, ipDefaultLimit int, ipTTL int) Limiter {
	limiter := Limiter{
		Storage:        storage,
		IpCooldown:     ipCooldown,
		IpDefaultLimit: ipDefaultLimit,
		IpTTL:          ipTTL,
	}

	return limiter
}

func (l *Limiter) HandleToken(tk string) (bool, error) {
	tokenConfig, err := l.Storage.GetTokenConfig(tk)

	if err != nil {
		return false, err
	}
	return l.check(tk, tokenConfig.ReqLimit, tokenConfig.TtlLimit, tokenConfig.Cooldown)
}

func (l *Limiter) HandleIP(ip string) (bool, error) {
	return l.check(ip, l.IpDefaultLimit, l.IpTTL, l.IpCooldown)
}

func (l *Limiter) check(term string, limit int, reqTTL int, cooldown int) (bool, error) {
	counter, err := l.Storage.GetCounter(term)

	if err != nil {
		return false, err
	}

	if counter == limit {
		l.Storage.RegisterBlock(term, cooldown)
		return false, nil
	}

	l.Storage.IncrementCounter(term, reqTTL)
	return true, nil
}
