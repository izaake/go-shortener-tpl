package urls

import (
	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/storage"
)

type sqlRepository struct {
	s storage.Storage
}

func NewSQLRepository(s storage.Storage) Repository {
	return &sqlRepository{s: s}
}

func (r sqlRepository) Save(user *models.User) error {
	//добавить транзакцию
	for _, u := range user.URLs {
		q := "INSERT INTO urls(user_id, short_url, original_url) values ($1, $2, $3) ON CONFLICT DO NOTHING"
		_, err := r.s.Exec(q, user.ID, u.ShortURL, u.FullURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r sqlRepository) FindOriginalURLByShortURL(url string) string {
	var shortURL string
	row := r.s.QueryRow("SELECT original_url FROM urls where short_url = $1 LIMIT 1", url)

	err := row.Scan(&shortURL)
	if err != nil {
		return shortURL
	}

	return shortURL
}

func (r sqlRepository) FindUrlsByUserID(userID string) []models.URL {
	urls := make([]models.URL, 0)

	rows, err := r.s.Query("SELECT short_url, original_url FROM urls where user_id = $1", userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var u models.URL
		err = rows.Scan(&u.ShortURL, &u.FullURL)
		if err != nil {
			return nil
		}

		urls = append(urls, u)
	}

	err = rows.Err()
	if err != nil {
		return nil
	}

	return urls
}

func (r sqlRepository) IsAvailable() bool {
	err := r.s.Ping()
	return err == nil
}
