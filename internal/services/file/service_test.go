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

	var expextedUsers []models.User
	expextedUsers = append(expextedUsers, models.User{Id: "123", URLs: expextedURLs}, models.User{Id: "12345", URLs: expextedURLs})

	for _, u := range expextedUsers {
		WriteToFile(filename, &u)
	}

	users, _ := ReadLines(filename)

	assert.Equal(t, expextedUsers, users)
}
