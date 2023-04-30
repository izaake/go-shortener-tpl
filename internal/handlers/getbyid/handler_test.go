package getbyid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCookie = "token=4a0b04b3-a2cb-4885-af15-9a342e817b00.f22b9af276e08f49c204b7a892cb5d211162255b0808dd891094c48a8f854e8a"

func TestHandler(t *testing.T) {
	// Получаем короткую ссылку для URL
	url := "https://practicum.yandex.ru"
	setResponse := testRequest(t, setshorturl.Handler, http.MethodPost, "/", strings.NewReader(url))

	assert.Equal(t, http.StatusCreated, setResponse.Code)
	assert.Equal(t, "/6bdb5b0e26a76e4dab7cd1a272caebc0", setResponse.Body.String())

	// По короткой ссылке получаем полный URL
	getResponse := testRequest(t, Handler, http.MethodGet, "/6bdb5b0e26a76e4dab7cd1a272caebc0", nil)

	assert.Equal(t, http.StatusTemporaryRedirect, getResponse.Code)
	assert.Equal(t, url, getResponse.Header().Get("location"))
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
			resp := testRequest(t, Handler, tt.method, tt.request, nil)
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
