package urls

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/izaake/go-shortener-tpl/internal/services/tokenutil"
)

// Handler — обработчик запроса на получение всех сохранённых ссылок юзера
func Handler(w http.ResponseWriter, r *http.Request) {
	splitUserToken := strings.Split(w.Header().Get("Set-Cookie"), "=")
	token := splitUserToken[1]
	userId, err := tokenutil.DecodeUserIdFromToken(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusNoContent)
		return
	}

	repo := urls.NewRepository()
	URLs := repo.FindUrlsByUserId(userId)

	if len(URLs) == 0 {
		http.Error(w, "no content", http.StatusNoContent)
		return
	}

	result, err := json.Marshal(URLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
