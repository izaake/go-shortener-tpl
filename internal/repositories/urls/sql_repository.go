package urls

import (
	"context"
	"time"

	"github.com/izaake/go-shortener-tpl/internal/models"
	"github.com/izaake/go-shortener-tpl/internal/storage"
)

type sqlRepository struct {
	s storage.Storage
}

func NewSQLRepository(s storage.Storage) Repository {
	return &sqlRepository{s: s}
}

func (r sqlRepository) Save(user *models.User, ignoreConflicts bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	q := `INSERT INTO urls(user_id, short_url, original_url) values ($1, $2, $3)`
	if ignoreConflicts {
		q = q + ` ON CONFLICT DO NOTHING`
	}

	tx, err := r.s.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range user.URLs {
		if err = r.s.Exec(ctx, stmt, user.ID, v.ShortURL, v.OriginalURL); err != nil {
			return err
		}
	}

	return tx.Commit()
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
		err = rows.Scan(&u.ShortURL, &u.OriginalURL)
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
