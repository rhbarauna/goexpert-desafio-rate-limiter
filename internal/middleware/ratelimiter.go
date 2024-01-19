package middleware

import (
	"net"
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

func getParsedIp(address string) string {
	parsedIP := net.ParseIP(address)

	if parsedIP.To4() == nil {
		return "127.0.0.1"
	}

	if parsedIP.To16().String() != "" {
		return parsedIP.To16().String()
	}

	return parsedIP.To4().String()
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := strings.Trim(r.Header.Get("API_KEY"), " ")
			address := r.Header.Get("X-Real-IP")

			if address == "" {
				address = r.RemoteAddr
			}

			ip := getParsedIp(address)

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
