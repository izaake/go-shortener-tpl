package batch

import (
	"encoding/json"
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

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := urlsRepository.NewMemoryRepository("")

	req := []NamedURL{
		{ID: "1", OriginalURL: "https://practicum.yandex.ru"},
		{ID: "2", OriginalURL: "https://www.test.ru"},
	}
	rawBody, err := json.Marshal(req)
	require.NoError(t, err)

	r, w := testRequest(t, http.MethodPost, "/api/shorten/batch", strings.NewReader(string(rawBody)))
	New(repo, "").Handle(w, r)

	res := []NamedURL{{}}
	err = json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t,
		[]NamedURL{
			{ID: "1", ShortURL: "/6bdb5b0e26a76e4dab7cd1a272caebc0"},
			{ID: "2", ShortURL: "/a1201c228ffa869cab2c19772afed576"},
		},
		res)
}

func testRequest(t *testing.T, method string, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	w.Header().Add("Set-Cookie", testCookie)

	r, err := http.NewRequest(method, path, body)
	require.NoError(t, err)

	return r, w
}
