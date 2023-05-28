package urls

import "github.com/izaake/go-shortener-tpl/internal/models"

// Repository Интерфейс для репозитория
type Repository interface {
	// Save сохраняет юзера с ссылками
	Save(user *models.User) error

	// FindOriginalURLByShortURL ищет полную ссылку по сокращённому варианту
	FindOriginalURLByShortURL(url string) string

	// FindUrlsByUserID ищет все сохранённые ссылки по юзеру
	FindUrlsByUserID(userID string) []models.URL

	// IsAvailable проверяет доступность хранилища
	IsAvailable() bool
}
