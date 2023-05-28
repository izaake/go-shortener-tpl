package urls

import (
	"log"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/services/file"
)

type fileRepository struct {
	fileName string
}

func NewFileRepository(fileName string) Repository {
	return &fileRepository{fileName: fileName}
}

func (r fileRepository) Save(user *models.User) error {
	filePath := r.fileName
	if filePath != "" {
		err := file.WriteToFile(filePath, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r fileRepository) FindOriginalURLByShortURL(url string) string {
	users, err := file.ReadLines(r.fileName)
	if err != nil {
		log.Print(err)
		return ""
	}

	for _, user := range users {
		for _, u := range user.URLs {
			if u.ShortURL == url {
				return u.FullURL
			}
		}
	}

	return ""
}

func (r fileRepository) FindUrlsByUserID(userID string) []models.URL {
	urls := make([]models.URL, 0)

	users, err := file.ReadLines(r.fileName)
	if err != nil {
		log.Print(err)
		return nil
	}
	for _, user := range users {
		if user.ID == userID {
			for _, url := range user.URLs {
				urls = append(urls, models.URL{ShortURL: "/" + url.ShortURL, FullURL: url.FullURL})
			}
		}
	}

	return urls
}

func (r fileRepository) IsAvailable() bool {
	return true
}
