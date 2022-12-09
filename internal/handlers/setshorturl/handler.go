package setshorturl

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
)

// Handler — обработчик запроса обмена полной ссылки на сокращённое значение.
func Handler(w http.ResponseWriter, r *http.Request) {
	u, err := validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortU := GetMD5Hash(u)

	repo := urls.NewRepository()
	repo.Save(models.URL{ShortURL: shortU, FullURL: u.String()})

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(repo.GetBaseURL() + "/" + shortU))
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
			log.Fatal(err)
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
