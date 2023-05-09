package getbyid

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izaake/go-shortener-tpl/internal/handlers/setshorturl"
	"github.com/izaake/go-shortener-tpl/internal/mock_storage"
	urlsRepository "github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCookie = "token=4a0b04b3-a2cb-4885-af15-9a342e817b00.f22b9af276e08f49c204b7a892cb5d211162255b0808dd891094c48a8f854e8a"

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_storage.NewMockStorage(ctrl)
	repo := urlsRepository.NewRepository(s)

	// Получаем короткую ссылку для URL
	url := "https://practicum.yandex.ru"

	r, w := testRequest(t, http.MethodPost, "/", strings.NewReader(url))
	setshorturl.New(repo).Handle(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "/6bdb5b0e26a76e4dab7cd1a272caebc0", w.Body.String())

	// По короткой ссылке получаем полный URL
	r1, w1 := testRequest(t, http.MethodGet, "/6bdb5b0e26a76e4dab7cd1a272caebc0", nil)
	New(repo).Handle(w1, r1)

	assert.Equal(t, http.StatusTemporaryRedirect, w1.Code)
	assert.Equal(t, url, w1.Header().Get("location"))
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_storage.NewMockStorage(ctrl)
			repo := urlsRepository.NewRepository(s)

			r, w := testRequest(t, tt.method, tt.request, nil)
			New(repo).Handle(w, r)

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
