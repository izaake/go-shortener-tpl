package setshorturl

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/caarlos0/env"
)

type Config struct {
	BaseURL string `env:"BASE_URL"`
}

var Str = map[string]string{}
var lock = sync.RWMutex{}

// Handler — обработчик запроса обмена полной ссылки на сокращённое значение.
func Handler(w http.ResponseWriter, r *http.Request) {
	u, err := validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortU := GetMD5Hash(u)

	lock.Lock()
	Str[shortU] = u.String()
	lock.Unlock()

	w.WriteHeader(http.StatusCreated)

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	baseURL := "http://localhost:8080"
	if cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}

	_, err = w.Write([]byte(baseURL + "/" + shortU))
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
	su, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	u, err := url.ParseRequestURI(string(su))
	if err != nil {
		return nil, err
	}

	return u, nil
}
