package shorten

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/izaake/go-shortener-tpl/internal/services/tokenutil"
)

// URLData содержит в себе полную версию ссылки
type URLData struct {
	URL string `json:"url,omitempty"`
}

// Response структура ответа на запрос
type Response struct {
	Result string `json:"result,omitempty"`
}

// Handler — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	u, err := validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	splitUserToken := strings.Split(w.Header().Get("Set-Cookie"), "=")
	token := splitUserToken[1]

	shortURL := GetMD5Hash(u.URL)
	repo := urls.NewRepository()

	userID, _ := tokenutil.DecodeUserIDFromToken(token)
	var uls []models.URL
	uls = append(uls, models.URL{FullURL: u.URL, ShortURL: shortURL})
	user := models.User{ID: userID, URLs: uls}

	err = repo.Save(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	res := Response{}
	res.Result = repo.GetBaseURL() + "/" + shortURL
	result, err := json.Marshal(res)
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

func GetMD5Hash(u string) string {
	hash := md5.Sum([]byte(u))
	return hex.EncodeToString(hash[:])
}

func validate(r *http.Request) (*URLData, error) {
	reader := r.Body
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		reader = gz
	}

	rawValue, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	u := URLData{}
	if err := json.Unmarshal(rawValue, &u); err != nil {
		return nil, err
	}

	_, err = url.ParseRequestURI(u.URL)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
