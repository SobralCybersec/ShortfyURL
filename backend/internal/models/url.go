package models

import "time"

type URL struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	Clicks      int64     `json:"clicks"`
}

type CreateURLRequest struct {
	URL string `json:"url"`
}

type CreateURLResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}
