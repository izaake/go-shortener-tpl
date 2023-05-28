package getbyid

import (
	"net/http"
	"strings"

	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
)

type Handler struct {
	repo urls.Repository
}

func New(
	repo urls.Repository,
) *Handler {
	return &Handler{
		repo: repo,
	}
}

// Handle — обработчик запроса поиска полной ссылки по сокращённому значению
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.Split(r.URL.Path, "/")[1]
	fullURL := h.repo.FindOriginalURLByShortURL(shortURL)

	if fullURL == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
