package models

// URLRequest is the request body for the URL shortening endpoint
type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

// URLResponse is the response body for the URL shortening endpoint
type URLResponse struct {
	ShortURL string `json:"short_url"`
}
