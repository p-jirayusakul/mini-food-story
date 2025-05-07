package http

type createResponse struct {
	ID int64 `json:"id"`
}

type createSessionResponse struct {
	URL string `json:"url"`
}
