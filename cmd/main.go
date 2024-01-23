package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/configs"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/limiter"
	ratelimiter "github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/middleware"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/storage/redis"
)

func main() {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Erro ao obter informações do arquivo.")
		return
	}
	goDir := filepath.Dir(currentFile)

	configs, err := configs.LoadConfig(goDir)
	if err != nil {
		panic(err)
	}

	// DEPS
	storage := redis.NewRedisStorage(configs.RedisHost, configs.RedisPort, configs.RedisPassword, configs.RedisDatabase)
	limiter := limiter.NewLimiter(storage, configs.Cooldown, configs.MaxRequests, configs.Ttl, configs.Tokens)
	rateLimiterMiddleware := ratelimiter.NewRateLimiter(limiter)
	// END DEPS

	// WEBSERVER
	router := chi.NewRouter()
	webServerPort := configs.WebServerPort

	router.Use(middleware.Logger)
	router.Use(rateLimiterMiddleware.Limit)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Iniciando o servidor web...")
	http.ListenAndServe(webServerPort, router)
	// END WEBSERVER
}
