package shorten

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCookie = "token=4a0b04b3-a2cb-4885-af15-9a342e817b00.f22b9af276e08f49c204b7a892cb5d211162255b0808dd891094c48a8f854e8a"

func TestHandler(t *testing.T) {
	u := URLData{
		URL: "https://practicum.yandex.ru",
	}
	rawBody, err := json.Marshal(u)
	require.NoError(t, err)

	resp := testRequest(t, Handler, http.MethodPost, "/api/shorten", strings.NewReader(string(rawBody)))

	res := Response{}
	err = json.Unmarshal(resp.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, "/6bdb5b0e26a76e4dab7cd1a272caebc0", res.Result)
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
			rawBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			resp := testRequest(t, Handler, http.MethodPost, "/api/shorten", strings.NewReader(string(rawBody)))
			assert.Equal(t, tt.want.statusCode, resp.Code)
			assert.Equal(t, tt.want.response, resp.Body.String())
		})
	}
}

func testRequest(t *testing.T, handler http.HandlerFunc, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	w.Header().Add("Set-Cookie", testCookie)

	r, err := http.NewRequest(method, path, body)
	require.NoError(t, err)

	handler.ServeHTTP(w, r)

	return w
}
