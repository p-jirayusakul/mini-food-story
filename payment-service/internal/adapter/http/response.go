package http

type createResponse struct {
	TransactionID string `json:"transactionID"`
}

type createSessionResponse struct {
	URL string `json:"url"`
}
