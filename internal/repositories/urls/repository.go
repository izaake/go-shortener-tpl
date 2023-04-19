package urls

import (
	"log"
	"sync"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/services/file"
)

// Repository Интерфейс для репозитория
type Repository interface {
	// Save сохраняет ссылку in memory
	Save(url *models.URL) error

	// Find ищет полную ссылку по сокращённому варианту
	Find(url string) string

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

var URLS = map[string]string{}
var BaseURL string
var FilePath string
var lock = sync.RWMutex{}

// NewRepository возвращает новый инстанс репозитория
func NewRepository() Repository {
	return &urlsRepository{}
}

// Save сохраняет ссылку
func (r urlsRepository) Save(url *models.URL) error {
	filePath := r.GetFilePath()
	if filePath != "" {
		err := file.WriteToFile(filePath, url)
		if err != nil {
			return err
		}
	}

	lock.Lock()
	URLS[url.ShortURL] = url.FullURL
	lock.Unlock()

	return nil
}

// Find ищет полную ссылку по сокращённому варианту
func (r urlsRepository) Find(url string) string {
	lock.RLock()
	u := URLS[url]
	lock.RUnlock()
	return u
}

// RestoreFromFile восстанавливает данные из файла в память
func (r urlsRepository) RestoreFromFile(filePath string) {
	if filePath != "" {
		urls, err := file.ReadLines(filePath)
		if err != nil {
			log.Print(err)
			return
		}
		for _, url := range urls {
			URLS[url.ShortURL] = url.FullURL
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
