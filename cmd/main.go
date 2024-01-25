package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	rateLimiterMiddleware, err := NewRateLimiter()

	if err != nil {
		panic(err)
	}

	// WEBSERVER
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(rateLimiterMiddleware.Limit)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Iniciando o servidor web...")
	http.ListenAndServe(GetWebServerPort(), router)
	// END WEBSERVER
}
