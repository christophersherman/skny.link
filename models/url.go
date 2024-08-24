package models

import "time"

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

type URLDataEntry struct {
	URL          string
	ShortURL     string
	CreatedAt    time.Time
	LastAccessed time.Time
	ViewCount    int
}
