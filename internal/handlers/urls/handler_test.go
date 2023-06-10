package urls

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/repositories/urls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name  string
		token string
		want  want
	}{
		{
			name:  "Получение сохранённых URL по известному юзеру",
			token: "4a0b04b3-a2cb-4885-af15-9a342e817b00.f22b9af276e08f49c204b7a892cb5d211162255b0808dd891094c48a8f854e8a",
			want: want{
				statusCode: http.StatusOK,
				response:   "[{\"short_url\":\"/bbb\",\"original_url\":\"aaa\"}]",
			},
		},
		{
			name:  "Получение сохранённых URL по юзеру без сохранённых ссылок",
			token: "123.123",
			want: want{
				statusCode: http.StatusNoContent,
				response:   "no content\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := urls.NewMemoryRepository("")

			var uls []models.URL
			uls = append(uls, models.URL{OriginalURL: "aaa", ShortURL: "bbb"})
			user := models.User{ID: "4a0b04b3-a2cb-4885-af15-9a342e817b00", URLs: uls}

			repo.Save(&user, true)

			r, w := testRequest(t, "/api/user/urls", tt.token)
			New(repo).Handle(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.response, w.Body.String())
		})
	}
}

func testRequest(t *testing.T, path string, token string) (*http.Request, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	w.Header().Add("Set-Cookie", "token="+token)

	r, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	return r, w
}
