package getbyid

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
)

var lock = sync.RWMutex{}

// Handler — обработчик запроса поиска полной ссылки по сокращённому значению.
func Handler(w http.ResponseWriter, r *http.Request) {
	shu := chi.URLParam(r, "id")

	lock.RLock()
	su := setshorturl.Str[shu]
	lock.RUnlock()

	if su == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", su)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
