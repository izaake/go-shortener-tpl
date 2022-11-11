package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/handlers"
)

func main() {
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", handler.Handler)
		r.Post("/", handler.Handler)
	})
	return r
}
