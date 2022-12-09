package file

import (
	"os"
	"testing"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	filename := "log"
	defer os.Remove(filename)

	var expextedURLs []models.URL
	expextedURLs = append(expextedURLs, models.URL{FullURL: "awdwd", ShortURL: "wedewdw"}, models.URL{FullURL: "1235", ShortURL: "12435"})

	for _, u := range expextedURLs {
		WriteToFile(filename, &u)
	}

	urls, _ := ReadLines(filename)

	assert.Equal(t, expextedURLs, urls)
}
