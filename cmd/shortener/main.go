package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/handlers/getbyid"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
	"github.com/izaake/go-shortener-tpl/internal/handlers/shorten"
)

func main() {
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
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
