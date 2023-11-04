package models

type ShortenedURL struct {
	URL          string `json:"url"`
	ShortenedURL string `json:"shortened_url"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}
