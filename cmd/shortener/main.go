package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/izaake/go-shortener-tpl/internal/handlers/getbyid"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
	"github.com/izaake/go-shortener-tpl/internal/handlers/shorten"
	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
}

const headerContentType = "Content-Type"

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	sAddr := flag.String("a", ":8080", "SERVER_ADDRESS")
	baseURL := flag.String("b", "http://localhost:8080", "BASE_URL")
	filePath := flag.String("f", "", "FILE_STORAGE_PATH")
	flag.Parse()

	if cfg.ServerAddress != "" {
		*sAddr = cfg.ServerAddress
	}
	if cfg.BaseURL != "" {
		*baseURL = cfg.BaseURL
	}
	if cfg.FilePath != "" {
		*filePath = cfg.FilePath
	}

	repo := urls.NewRepository()
	// Восстанавливаем сохранённые url из файла
	repo.RestoreFromFile(*filePath)
	repo.SaveBaseURL(*baseURL)
	repo.SaveFilePath(*filePath)

	r := NewRouter()
	log.Fatal(http.ListenAndServe(*sAddr, r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(commonMiddleware)
	r.Use(middleware.Compress(5))

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", getbyid.Handler)
		r.Post("/", setshorturl.Handler)
		r.Post("/api/shorten", shorten.Handler)
	})
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerContentType) == "application/json" {
			w.Header().Add(headerContentType, "application/json")
		}
		next.ServeHTTP(w, r)
	})
}
