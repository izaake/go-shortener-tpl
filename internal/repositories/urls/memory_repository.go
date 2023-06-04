package urls

import (
	"sync"

	"github.com/izaake/go-shortener-tpl/internal/models"
)

type memoryRepository struct {
	Users   map[string]map[string]string
	BaseURL string
}

var (
	lock = sync.RWMutex{}
)

func NewMemoryRepository(baseURL string) Repository {
	Users := make(map[string]map[string]string)
	return &memoryRepository{Users: Users, BaseURL: baseURL}
}

func (r memoryRepository) Save(user *models.User) error {
	if r.Users[user.ID] == nil {
		r.Users[user.ID] = map[string]string{}
	}

	for _, url := range user.URLs {
		lock.Lock()
		r.Users[user.ID][url.ShortURL] = url.OriginalURL
		lock.Unlock()
	}
	return nil
}

func (r memoryRepository) FindOriginalURLByShortURL(url string) string {
	var u string
	for _, user := range r.Users {
		lock.RLock()
		if user[url] != "" {
			u = user[url]
		}
		lock.RUnlock()
	}
	return u
}

func (r memoryRepository) FindUrlsByUserID(userID string) []models.URL {
	urls := make([]models.URL, 0)

	lock.RLock()
	for k, v := range r.Users[userID] {
		urls = append(urls, models.URL{ShortURL: r.BaseURL + "/" + k, OriginalURL: v})
	}
	lock.RUnlock()

	return urls
}

func (r memoryRepository) IsAvailable() bool {
	return true
}
