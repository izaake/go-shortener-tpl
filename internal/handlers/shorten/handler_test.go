package shorten

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)

	u := URLData{
		URL: "https://practicum.yandex.ru",
	}
	rawBody, err := json.Marshal(u)
	require.NoError(t, err)

	statusCode, respBody := testRequest(t, ts, "/api/shorten", strings.NewReader(string(rawBody)))

	res := Response{}
	err = json.Unmarshal([]byte(respBody), &res)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "http://localhost:8080/6bdb5b0e26a76e4dab7cd1a272caebc0", res.Result)
}

func TestHandlerNegative(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name        string
		requestBody URLData
		want        want
	}{
		{
			name: "POST - Переданный URL в запросе невалиден",
			requestBody: URLData{
				URL: "123",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "parse \"123\": invalid URI for request\n",
			},
		},
		{
			name:        "POST - URL в запросе не передан",
			requestBody: URLData{},
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

			rawBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			statusCode, body := testRequest(t, ts, "/api/shorten", strings.NewReader(string(rawBody)))
			assert.Equal(t, tt.want.statusCode, statusCode)
			assert.Equal(t, tt.want.response, body)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, path string, body io.Reader) (int, string) {
	req, err := http.NewRequest(http.MethodPost, ts.URL+path, body)
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
		r.Post("/api/shorten", Handler)
	})
	return r
}
