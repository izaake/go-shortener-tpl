package getbyid

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
)

// Handler — обработчик запроса поиска полной ссылки по сокращённому значению.
func Handler(w http.ResponseWriter, r *http.Request) {
	shu := chi.URLParam(r, "id")

	repo := urls.NewRepository()
	su := repo.Find(shu)

	if su == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", su)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
