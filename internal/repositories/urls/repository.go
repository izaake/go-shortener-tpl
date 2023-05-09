package urls

import (
	"log"
	"sync"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/services/file"
	"github.com/izaake/go-shortener-tpl/internal/storage"
)

// Repository Интерфейс для репозитория
// todo разделить работу с файлом и с ссылками
type Repository interface {
	// Save сохраняет ссылку in memory
	Save(user *models.User) error

	// FindOriginalURLByShortURL ищет полную ссылку по сокращённому варианту
	FindOriginalURLByShortURL(url string) string

	// FindUrlsByUserID ищет все сохранённые ссылки по юзеру
	FindUrlsByUserID(userID string) []models.URL

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

	PingDB() error
}

var (
	Users    = map[string]map[string]string{}
	BaseURL  string
	FilePath string
	lock     = sync.RWMutex{}
)

type urlsRepository struct {
	s storage.Storage
}

// NewRepository возвращает новый инстанс репозитория
func NewRepository(s storage.Storage) Repository {
	return &urlsRepository{s}
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

	if Users[user.ID] == nil {
		Users[user.ID] = map[string]string{}
	}

	for _, url := range user.URLs {
		lock.Lock()
		Users[user.ID][url.ShortURL] = url.FullURL
		lock.Unlock()
	}

	return nil
}

// FindOriginalURLByShortURL ищет полную ссылку по сокращённому варианту
func (r urlsRepository) FindOriginalURLByShortURL(url string) string {
	var u string
	for _, user := range Users {
		lock.RLock()
		if user[url] != "" {
			u = user[url]
		}
		lock.RUnlock()
	}
	return u
}

// FindUrlsByUserID все сохранённые ссылки по юзеру
func (r urlsRepository) FindUrlsByUserID(userID string) []models.URL {
	urls := make([]models.URL, 0)

	lock.RLock()
	for k, v := range Users[userID] {
		urls = append(urls, models.URL{ShortURL: r.GetBaseURL() + "/" + k, FullURL: v})
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
				if Users[user.ID] == nil {
					Users[user.ID] = map[string]string{}
				}
				Users[user.ID][url.ShortURL] = url.FullURL
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

func (r urlsRepository) PingDB() error {
	return r.s.Ping()
}
