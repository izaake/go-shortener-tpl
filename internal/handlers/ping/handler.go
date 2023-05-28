package ping

import (
	"net/http"

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

// Handle — обработчик запроса на получение всех сохранённых ссылок юзера
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if !h.repo.IsAvailable() {
		http.Error(w, "db ping error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
