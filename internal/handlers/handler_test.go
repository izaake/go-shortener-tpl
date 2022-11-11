package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	// Получаем короткую ссылку для URL
	url := "https://practicum.yandex.ru"
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url))

	w := httptest.NewRecorder()
	h := http.HandlerFunc(Handler)
	h(w, request)

	result := w.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	resBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	err = result.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8080/6bdb5b0e26a76e4dab7cd1a272caebc0", string(resBody))

	// По короткой ссылке получаем полный URL
	requestGet := httptest.NewRequest(http.MethodGet, "/6bdb5b0e26a76e4dab7cd1a272caebc0", nil)
	w = httptest.NewRecorder()
	h(w, requestGet)

	result = w.Result()
	assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
	assert.Equal(t, url, result.Header.Get("location"))
	err = result.Body.Close()
	require.NoError(t, err)
}

func TestHandlerGetNegative(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name    string
		request string
		schema  string
		want    want
	}{
		{
			name:    "GET - Отсутствует ID сокращённого URL в запросе",
			request: "/",
			schema:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "ID is missing\n",
			},
		},
		{
			name:    "GET - URL по заданному ID не найден",
			request: "/123",
			schema:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "not found\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.schema, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(Handler)
			h.ServeHTTP(w, request)

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
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
		schema      string
		want        want
	}{
		{
			name:        "POST - Переданный URL в запросе невалиден",
			requestBody: "123",
			schema:      http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "parse \"123\": invalid URI for request\n",
			},
		},
		{
			name:        "POST - URL в запросе не передан",
			requestBody: "",
			schema:      http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "parse \"\": empty url\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.schema, "/", strings.NewReader(tt.requestBody))

			w := httptest.NewRecorder()
			h := http.HandlerFunc(Handler)
			h.ServeHTTP(w, request)

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
