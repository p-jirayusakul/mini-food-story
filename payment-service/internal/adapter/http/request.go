package http

type Payment struct {
	OrderID string  `json:"orderID" validate:"required"`
	Method  string  `json:"methodID" validate:"required"`
	Note    *string `json:"note"`
}

type CallbackPayment struct {
	TransactionID string `json:"transactionID" validate:"required"`
}
