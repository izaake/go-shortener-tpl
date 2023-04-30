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

func Test_urlsRepository_FindOriginalUrlByUserId(t *testing.T) {
	expectedURL := "awdwd"
	shortUrl := "wedewdw"
	userId := uuid.New().String()

	URL := make([]models.URL, 0)
	URL = append(URL, models.URL{FullURL: expectedURL, ShortURL: shortUrl})

	user := models.User{Id: userId, URLs: URL}
	repo := NewRepository()
	repo.Save(&user)
	actualURL := repo.FindOriginalUrlByShortUrl(shortUrl)

	assert.Equal(t, expectedURL, actualURL)
}

func Test_urlsRepository_RestoreFromFile(t *testing.T) {
	filename := "u.log"
	defer os.Remove(filename)

	short := "wedewdw"
	orig := "awdwd"

	var expextedURLs []models.URL
	expextedURLs = append(expextedURLs, models.URL{FullURL: orig, ShortURL: short})
	expectedUser := models.User{Id: "111", URLs: expextedURLs}

	repo := NewRepository()
	file.WriteToFile("u.log", &expectedUser)
	repo.RestoreFromFile("u.log")
	actualURL := repo.FindOriginalUrlByShortUrl(short)

	assert.Equal(t, orig, actualURL)
}

func Test_urlsRepository_FindUrlsByUserId(t *testing.T) {
	userId := uuid.New().String()
	expectedURLs := []models.URL{
		{FullURL: "123", ShortURL: "321"},
		{FullURL: "qwe", ShortURL: "ewq"},
	}

	user := models.User{Id: userId, URLs: expectedURLs}
	repo := NewRepository()
	repo.Save(&user)

	actualURLs := repo.FindUrlsByUserId(userId)

	assert.Equal(t, true, reflect.DeepEqual(expectedURLs, actualURLs))
}
