package urls

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/services/file"
	"github.com/stretchr/testify/assert"
)

func Test_urlsRepository_FindOriginalUrlByUserID(t *testing.T) {
	expectedURL := "awdwd"
	shortURL := "wedewdw"
	userID := uuid.New().String()

	URL := make([]models.URL, 0)
	URL = append(URL, models.URL{FullURL: expectedURL, ShortURL: shortURL})

	user := models.User{ID: userID, URLs: URL}
	repo := NewRepository()
	repo.Save(&user)
	actualURL := repo.FindOriginalURLByShortURL(shortURL)

	assert.Equal(t, expectedURL, actualURL)
}

func Test_urlsRepository_RestoreFromFile(t *testing.T) {
	filename := "u.log"
	defer os.Remove(filename)

	short := "wedewdw"
	orig := "awdwd"

	var expextedURLs []models.URL
	expextedURLs = append(expextedURLs, models.URL{FullURL: orig, ShortURL: short})
	expectedUser := models.User{ID: "111", URLs: expextedURLs}

	repo := NewRepository()
	file.WriteToFile("u.log", &expectedUser)
	repo.RestoreFromFile("u.log")
	actualURL := repo.FindOriginalURLByShortURL(short)

	assert.Equal(t, orig, actualURL)
}

func Test_urlsRepository_FindUrlsByUserID(t *testing.T) {
	userID := uuid.New().String()
	urls := []models.URL{
		{FullURL: "123", ShortURL: "321"},
		{FullURL: "qwe", ShortURL: "ewq"},
	}
	expectedURLs := []models.URL{
		{FullURL: "123", ShortURL: "/321"},
		{FullURL: "qwe", ShortURL: "/ewq"},
	}

	user := models.User{ID: userID, URLs: urls}
	repo := NewRepository()
	repo.Save(&user)

	actualURLs := repo.FindUrlsByUserID(userID)

	assert.Equal(t, true, reflect.DeepEqual(expectedURLs, actualURLs))
}
