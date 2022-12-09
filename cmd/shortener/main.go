package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/handlers/getbyid"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
	"github.com/izaake/go-shortener-tpl/internal/handlers/shorten"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	sAddr := ":8080"
	if cfg.ServerAddress != "" {
		sAddr = cfg.ServerAddress
	}

	r := NewRouter()
	log.Fatal(http.ListenAndServe(sAddr, r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(commonMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", getbyid.Handler)
		r.Post("/", setshorturl.Handler)
		r.Post("/api/shorten", shorten.Handler)
	})
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
