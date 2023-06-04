package models

type URL struct {
	ShortURL    string `json:"short_url"`    // Результирующий сокращённый URL
	OriginalURL string `json:"original_url"` // URL для сокращения
}
