package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/handlers/get_by_id"
	"github.com/izaake/go-shortener-tpl/internal/handlers/set_short_url"
)

func main() {
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", get_by_id.Handler)
		r.Post("/", set_short_url.Handler)
	})
	return r
}
