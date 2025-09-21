package http

type createResponse struct {
	TransactionID string `json:"transactionID" example:"e88ac32b-c940-4102-a51f-c7a2e7ed6622"`
}

type LastStatusCodeResponse struct {
	Code string `json:"code" example:"SUCCESS"`
}
