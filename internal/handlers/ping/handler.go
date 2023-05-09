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
	err := h.repo.PingDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
