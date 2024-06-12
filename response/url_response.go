package response

type UrlResponse struct {
	OriginalURL string `json:'original_url'`
	ShortURL    string `json:'short_url'`
	Hits        int    `json:'hits'`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
