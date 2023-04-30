package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/izaake/go-shortener-tpl/internal/handlers/getbyid"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
	"github.com/izaake/go-shortener-tpl/internal/handlers/shorten"
	"github.com/izaake/go-shortener-tpl/internal/handlers/urls"
	urlsRepository "github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/izaake/go-shortener-tpl/internal/services/tokenutil"
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

	repo := urlsRepository.NewRepository()
	// Восстанавливаем сохранённые url по сохранённым юзерам из файла
	repo.RestoreFromFile(*filePath)
	repo.SaveBaseURL(*baseURL)
	repo.SaveFilePath(*filePath)

	r := NewRouter()
	log.Fatal(http.ListenAndServe(*sAddr, r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Use(commonMiddleware)
	r.Use(middleware.Compress(5))

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", getbyid.Handler)
		r.Post("/", setshorturl.Handler)
		r.Post("/api/shorten", shorten.Handler)
		r.Get("/api/user/urls", urls.Handler)
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cfg Config
		err := env.Parse(&cfg)
		if err != nil {
			log.Fatal(err)
		}

		var token string
		t, err := r.Cookie("token") // строка в формате "token=<user_id>.<sign>"
		if t != nil {
			splitToken := strings.Split(t.String(), "=")
			token = splitToken[1]
			_, err = tokenutil.DecodeUserIdFromToken(token)
		}

		if err != nil || tokenutil.IsTokenValid(token) == false {
			userId := tokenutil.GenerateUserId()
			token = tokenutil.GenerateTokenForUser(userId)
		}

		http.SetCookie(w, &http.Cookie{Name: "token", Value: token})

		next.ServeHTTP(w, r)
	})
}
