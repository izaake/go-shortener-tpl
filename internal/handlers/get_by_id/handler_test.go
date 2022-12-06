package get_by_id

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/izaake/go-shortener-tpl/internal/handlers/set_short_url"
	"github.com/izaake/go-shortener-tpl/internal/handlers/shorten"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	// Получаем короткую ссылку для URL
	url := "https://practicum.yandex.ru"
	r := NewRouter()
	ts := httptest.NewServer(r)

	statusCode, body := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(url))
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "http://localhost:8080/6bdb5b0e26a76e4dab7cd1a272caebc0", body)

	// По короткой ссылке получаем полный URL
	statusCode, _ = testRequest(t, ts, http.MethodGet, "/6bdb5b0e26a76e4dab7cd1a272caebc0", nil)
	assert.Equal(t, http.StatusOK, statusCode)
	ts.Close()
}

func TestHandlerGetNegative(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "GET - URL по заданному ID не найден",
			request: "/123",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "not found\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, body := testRequest(t, ts, tt.method, tt.request, nil)
			assert.Equal(t, tt.want.statusCode, statusCode)
			assert.Equal(t, tt.want.response, body)
		})
	}
}

func TestHandlerPostNegative(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name        string
		requestBody string
		method      string
		want        want
	}{
		{
			name:        "POST - Переданный URL в запросе невалиден",
			requestBody: "123",
			method:      http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "parse \"123\": invalid URI for request\n",
			},
		},
		{
			name:        "POST - URL в запросе не передан",
			requestBody: "",
			method:      http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "parse \"\": empty url\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, body := testRequest(t, ts, tt.method, "/", strings.NewReader(tt.requestBody))
			assert.Equal(t, tt.want.statusCode, statusCode)
			assert.Equal(t, tt.want.response, body)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method string, path string, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", Handler)
		r.Post("/", set_short_url.Handler)
		r.Post("/api/shorten", shorten.Handler)
	})
	return r
}
