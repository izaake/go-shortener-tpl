package getbyid

import (
	"net/http"
	"strings"

	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/izaake/go-shortener-tpl/internal/services/tokenutil"
)

// Handler — обработчик запроса поиска полной ссылки по сокращённому значению
func Handler(w http.ResponseWriter, r *http.Request) {
	splitUserToken := strings.Split(w.Header().Get("Set-Cookie"), "=")
	token := splitUserToken[1]
	userId, _ := tokenutil.DecodeUserIdFromToken(token)

	shortURL := strings.Split(r.URL.Path, "/")[1]

	repo := urls.NewRepository()
	fullURL := repo.FindOriginalUrlByUserId(shortURL, userId)

	if fullURL == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
