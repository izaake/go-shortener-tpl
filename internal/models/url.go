package models

type URL struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}
