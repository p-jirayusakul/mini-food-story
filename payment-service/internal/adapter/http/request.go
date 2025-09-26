package http

type Payment struct {
	OrderID string  `json:"orderID" validate:"required" example:"1921144250070732800"`
	Method  string  `json:"methodID" validate:"required" example:"1923732004537372675"`
	Note    *string `json:"note" example:"lorem ipsum"`
}

type CallbackPayment struct {
	TransactionID string `json:"transactionID" validate:"required" example:"e88ac32b-c940-4102-a51f-c7a2e7ed6622"`
	Status        string `json:"status" validate:"required" example:"SUCCESS"`
}
