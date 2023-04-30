package urls

import (
	"log"
	"sync"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/services/file"
)

// Repository Интерфейс для репозитория
// todo разделить работу с файлом и с ссылками
type Repository interface {
	// Save сохраняет ссылку in memory
	Save(user *models.User) error

	// FindOriginalUrlByUserId ищет полную ссылку по сокращённому варианту по юзеру
	FindOriginalUrlByUserId(url string, userId string) string

	// FindUrlsByUserId ищет все сохранённые ссылки по юзеру
	FindUrlsByUserId(userId string) []models.URL

	// RestoreFromFile восстанавливает данные из файла в память
	RestoreFromFile(filePath string)

	// SaveBaseURL сохраняет базовую ссылку
	SaveBaseURL(baseURL string)

	// GetBaseURL получить базовую ссылку
	GetBaseURL() string

	// SaveFilePath сохраняет путь до файла urls
	SaveFilePath(filePath string)

	// GetFilePath получить путь до файла urls
	GetFilePath() string
}

type urlsRepository struct{}

var Users = map[string]map[string]string{}
var BaseURL string
var FilePath string
var lock = sync.RWMutex{}

// NewRepository возвращает новый инстанс репозитория
func NewRepository() Repository {
	return &urlsRepository{}
}

// Save сохраняет модель юзера со ссылками
func (r urlsRepository) Save(user *models.User) error {
	filePath := r.GetFilePath()
	if filePath != "" {
		err := file.WriteToFile(filePath, user)
		if err != nil {
			return err
		}
	}

	if Users[user.Id] == nil {
		Users[user.Id] = map[string]string{}
	}

	for _, url := range user.URLs {
		lock.Lock()
		Users[user.Id][url.ShortURL] = url.FullURL
		lock.Unlock()
	}

	return nil
}

// FindOriginalUrlByUserId ищет полную ссылку по сокращённому варианту по юзеру
func (r urlsRepository) FindOriginalUrlByUserId(url string, userId string) string {
	lock.RLock()
	u := Users[userId][url]
	lock.RUnlock()
	return u
}

// FindUrlsByUserId все сохранённые ссылки по юзеру
func (r urlsRepository) FindUrlsByUserId(userId string) []models.URL {
	urls := make([]models.URL, 0)

	lock.RLock()
	for k, v := range Users[userId] {
		urls = append(urls, models.URL{ShortURL: k, FullURL: v})
	}
	lock.RUnlock()

	return urls
}

// RestoreFromFile восстанавливает данные из файла в память
func (r urlsRepository) RestoreFromFile(filePath string) {
	if filePath != "" {
		users, err := file.ReadLines(filePath)
		if err != nil {
			log.Print(err)
			return
		}
		for _, user := range users {
			for _, url := range user.URLs {
				if Users[user.Id] == nil {
					Users[user.Id] = map[string]string{}
				}
				Users[user.Id][url.ShortURL] = url.FullURL
			}
		}
	}
}

func (r urlsRepository) SaveBaseURL(baseURL string) {
	BaseURL = baseURL
}

func (r urlsRepository) GetBaseURL() string {
	return BaseURL
}

func (r urlsRepository) SaveFilePath(filePath string) {
	FilePath = filePath
}

func (r urlsRepository) GetFilePath() string {
	return FilePath
}
