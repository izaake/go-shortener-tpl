package setshorturl

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/izaake/go-shortener-tpl/internal/services/tokenutil"
)

type Handler struct {
	repo    urls.Repository
	baseURL string
}

func New(
	repo urls.Repository,
	baseURL string,
) *Handler {
	return &Handler{
		repo:    repo,
		baseURL: baseURL,
	}
}

// Handle — обработчик запроса обмена полной ссылки на сокращённое значение.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	fullURL, err := validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie := w.Header().Get("Set-Cookie")
	if cookie == "" {
		http.Error(w, "empty cookie header", http.StatusBadRequest)
		return
	}
	splitUserToken := strings.Split(cookie, "=")
	token := splitUserToken[1]

	userID, _ := tokenutil.DecodeUserIDFromToken(token)
	shortURL := GetMD5Hash(fullURL)

	var uls []models.URL
	uls = append(uls, models.URL{OriginalURL: fullURL.String(), ShortURL: shortURL})
	user := models.User{ID: userID, URLs: uls}

	err = h.repo.Save(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(h.baseURL + "/" + shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetMD5Hash(u *url.URL) string {
	hash := md5.Sum([]byte(u.String()))
	return hex.EncodeToString(hash[:])
}

func validate(r *http.Request) (*url.URL, error) {
	reader := r.Body
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		reader = gz
	}

	su, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	u, err := url.ParseRequestURI(string(su))
	if err != nil {
		return nil, err
	}

	return u, nil
}
