package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/limiter"
)

type RateLimiter struct {
	Limiter limiter.Limiter
}

func NewRateLimiter(limiter limiter.Limiter) *RateLimiter {
	return &RateLimiter{
		Limiter: limiter,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var moveOn bool
			var err error

			token := strings.Trim(r.Header.Get("API_KEY"), " ")

			if token != "" {
				moveOn, err = rl.Limiter.HandleToken(token)
			}

			if token == "" || err != nil {
				ip := ""
				moveOn, err = rl.Limiter.HandleIP(ip)
			}

			if moveOn {
				log.Printf("Request bloqueada %v\n", r)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err != nil {
				log.Printf("Erro no limiter %v\n", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
}
