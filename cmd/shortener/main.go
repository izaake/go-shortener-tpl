package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/izaake/go-shortener-tpl/internal/database"
	"github.com/izaake/go-shortener-tpl/internal/handlers/batch"
	"github.com/izaake/go-shortener-tpl/internal/handlers/getbyid"
	"github.com/izaake/go-shortener-tpl/internal/handlers/ping"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
	"github.com/izaake/go-shortener-tpl/internal/handlers/shorten"
	"github.com/izaake/go-shortener-tpl/internal/handlers/urls"
	urlsRepository "github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/izaake/go-shortener-tpl/internal/services/tokenutil"
	"github.com/izaake/go-shortener-tpl/internal/storage"
)

var (
	sAddr        = flag.String("a", ":8080", "SERVER_ADDRESS")
	baseURL      = flag.String("b", "http://localhost:8080", "BASE_URL")
	filePath     = flag.String("f", "", "FILE_STORAGE_PATH")
	dbConnection = flag.String("d", "postgres://eaizaak@localhost:5432/test", "DATABASE_DSN")
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	DBConnection  string `env:"DATABASE_DSN"`
}

const headerContentType = "Content-Type"

func main() {
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddress != "" {
		*sAddr = cfg.ServerAddress
	}
	if cfg.BaseURL != "" {
		*baseURL = cfg.BaseURL
	}
	if cfg.FilePath != "" {
		*filePath = cfg.FilePath
	}
	if cfg.DBConnection != "" {
		*dbConnection = cfg.DBConnection
	}

	repo := urlsRepository.NewMemoryRepository(*baseURL)
	if filePath != nil && *filePath != "" {
		repo = urlsRepository.NewFileRepository(*filePath)
	}
	if dbConnection != nil && *dbConnection != "" {
		db, err := database.NewDB(dbConnection)
		if err == nil {
			defer db.Close()

			st := &storage.SQLStorage{DB: db}
			repo = urlsRepository.NewSQLRepository(st)
		}
	}

	r := NewRouter(repo, *baseURL)
	log.Fatal(http.ListenAndServe(*sAddr, r))
}

func NewRouter(repo urlsRepository.Repository, baseURL string) chi.Router {
	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Use(commonMiddleware)
	r.Use(middleware.Compress(5))

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", getbyid.New(repo).Handle)
		r.Post("/", setshorturl.New(repo, baseURL).Handle)
		r.Post("/api/shorten", shorten.New(repo, baseURL).Handle)
		r.Post("/api/shorten/batch", batch.New(repo, baseURL).Handle)
		r.Get("/api/user/urls", urls.New(repo).Handle)
		r.Get("/ping", ping.New(repo).Handle)
	})
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerContentType) == "application/json" || r.URL.Path == "/api/user/urls" {
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
			_, err = tokenutil.DecodeUserIDFromToken(token)
		}

		if err != nil || !tokenutil.IsTokenValid(token) {
			userID := tokenutil.GenerateUserID()
			token = tokenutil.GenerateTokenForUser(userID)
		}

		http.SetCookie(w, &http.Cookie{Name: "token", Value: token})

		next.ServeHTTP(w, r)
	})
}
