package get_by_id

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/handlers/set_short_url"
)

var lock = sync.RWMutex{}

// Handler — обработчик запроса поиска полной ссылки по сокращённому значению.
func Handler(w http.ResponseWriter, r *http.Request) {
	shu := chi.URLParam(r, "id")

	lock.RLock()
	su := set_short_url.Str[shu]
	lock.RUnlock()

	if su == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", su)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}
