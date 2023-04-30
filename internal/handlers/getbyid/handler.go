package getbyid

import (
	"net/http"
	"strings"

	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
)

// Handler — обработчик запроса поиска полной ссылки по сокращённому значению
func Handler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.Split(r.URL.Path, "/")[1]

	repo := urls.NewRepository()
	fullURL := repo.FindOriginalURLByShortURL(shortURL)

	if fullURL == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
