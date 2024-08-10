package models

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}
