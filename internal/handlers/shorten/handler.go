package shorten

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/caarlos0/env"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
)

var lock = sync.RWMutex{}

// URLData содержит в себе полную версию ссылки
type URLData struct {
	URL string `json:"url,omitempty"`
}

// Response структура ответа на запрос
type Response struct {
	Result string `json:"result,omitempty"`
}

type Config struct {
	BaseURL string `env:"BASE_URL"`
}

// Handler — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	u, err := validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortU := GetMD5Hash(u.URL)

	lock.Lock()
	setshorturl.Str[shortU] = u.URL
	lock.Unlock()

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	baseURL := "http://localhost:8080"
	if cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}

	res := Response{}
	res.Result = baseURL + "/" + shortU
	result, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
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
	rawValue, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

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
