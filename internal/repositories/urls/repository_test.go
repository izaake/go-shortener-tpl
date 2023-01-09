package urls

import (
	"os"
	"testing"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/services/file"
	"github.com/stretchr/testify/assert"
)

func TestUrlsRepository_Find(t *testing.T) {
	expextedURL := models.URL{FullURL: "awdwd", ShortURL: "wedewdw"}

	repo := NewRepository()
	repo.Save(expextedURL)
	actualURL := repo.Find(expextedURL.ShortURL)

	assert.Equal(t, expextedURL.FullURL, actualURL)
}

func TestUrlsRepository_RestoreFromFile(t *testing.T) {
	filename := "u.log"
	defer os.Remove(filename)

	expextedURL := models.URL{FullURL: "awdwd", ShortURL: "wedewdw"}

	repo := NewRepository()
	file.WriteToFile("u.log", &expextedURL)
	repo.RestoreFromFile("u.log")
	actualURL := repo.Find(expextedURL.ShortURL)

	assert.Equal(t, expextedURL.FullURL, actualURL)
}
