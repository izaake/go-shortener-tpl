package handler

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
)

var str = map[string]string{}

// Handler — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		su, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := url.ParseRequestURI(string(su))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shu := GetMD5Hash(u)
		str[shu] = u.String()

		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte("http://localhost:8080/" + shu))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodGet:
		shu := r.URL.Path[1:]
		if shu == "" {
			http.Error(w, "ID is missing", http.StatusBadRequest)
			return
		}
		su := str[shu]

		if su == "" {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("location", su)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	default:
		http.Error(w, "Only GET/POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func GetMD5Hash(u *url.URL) string {
	hash := md5.Sum([]byte(u.String()))
	return hex.EncodeToString(hash[:])
}
