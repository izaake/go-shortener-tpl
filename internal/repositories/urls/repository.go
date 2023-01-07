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
	Save(filePath string, url models.URL)

	// Find ищет полную ссылку по сокращённому варианту
	Find(url string) string

	// RestoreFromFile восстанавливает данные из файла в память
	RestoreFromFile(filePath string)
}

type urlsRepository struct{}
type Config struct {
	FilePath string `env:"FILE_STORAGE_PATH"`
}

var URLS = map[string]string{}
var lock = sync.RWMutex{}

// NewRepository возвращает новый инстанс репозитория
func NewRepository() Repository {
	return &urlsRepository{}
}

// Save сохраняет ссылку
func (r urlsRepository) Save(filePath string, url models.URL) {
	if filePath != "" {
		file.WriteToFile(filePath, &url)
	}

	lock.Lock()
	URLS[url.ShortURL] = url.FullURL
	lock.Unlock()
}

// Find ищет полную ссылку по сокращённому варианту
func (r urlsRepository) Find(url string) string {
	lock.RLock()
	u := URLS[url]
	lock.RUnlock()
	return u
}

// RestoreFromFile сохраняет ссылку
func (r urlsRepository) RestoreFromFile(filePath string) {
	if filePath != "" {
		urls, err := file.ReadLines(filePath)
		if err != nil {
			log.Fatalf("readLines: %s", err)
		}
		for _, url := range urls {
			URLS[url.ShortURL] = url.FullURL
		}
	}
}
