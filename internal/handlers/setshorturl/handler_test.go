package setshorturl

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	urlsRepository "github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCookie = "token=4a0b04b3-a2cb-4885-af15-9a342e817b00.f22b9af276e08f49c204b7a892cb5d211162255b0808dd891094c48a8f854e8a"

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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := urlsRepository.NewMemoryRepository("")

			r, w := testRequest(t, tt.method, "/", strings.NewReader(tt.requestBody))
			New(repo, "").Handle(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.response, w.Body.String())
		})
	}
}

func testRequest(t *testing.T, method string, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	w.Header().Add("Set-Cookie", testCookie)

	r, err := http.NewRequest(method, path, body)
	require.NoError(t, err)

	return r, w
}
