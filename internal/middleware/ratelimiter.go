package middleware

import (
	"net/http"
	"strings"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/limiter"
)

type RateLimiter struct {
	limiter limiter.Limiter
}

func NewRateLimiter(limiter limiter.Limiter) RateLimiter {
	return RateLimiter{limiter: limiter}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := strings.Trim(r.Header.Get("API_KEY"), " ")
			ip := r.RemoteAddr

			err := rl.limiter.Limit(ip, token)

			if err == limiter.ErrLimitedAccess {
				http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame.", http.StatusTooManyRequests)
				return
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
}
