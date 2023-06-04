package batch

import (
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

type BatchIn struct {
	Batch []NamedURL
}

type NamedURL struct {
	ID          string `json:"correlation_id"`         // Строковый идентификатор
	OriginalURL string `json:"original_url,omitempty"` // URL для сокращения
	ShortURL    string `json:"short_url,omitempty"`    // Результирующий сокращённый URL
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	u, err := validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	splitUserToken := strings.Split(w.Header().Get("Set-Cookie"), "=")
	token := splitUserToken[1]

	var uls []models.URL
	out := make([]NamedURL, 0)
	for _, v := range u {
		su := GetMD5Hash(v.OriginalURL)
		uls = append(uls, models.URL{OriginalURL: v.OriginalURL, ShortURL: su})

		out = append(out, NamedURL{ID: v.ID, ShortURL: h.baseURL + "/" + su})
	}

	userID, _ := tokenutil.DecodeUserIDFromToken(token)
	user := models.User{ID: userID, URLs: uls}

	err = h.repo.Save(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	result, err := json.Marshal(out)
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

func validate(r *http.Request) ([]NamedURL, error) {
	reader := r.Body
	rawValue, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	u := make([]NamedURL, 0)
	if err := json.Unmarshal(rawValue, &u); err != nil {
		return nil, err
	}

	for _, v := range u {
		_, err = url.ParseRequestURI(v.OriginalURL)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}
