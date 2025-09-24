package domain

type Payment struct {
	ID      int64
	OrderID int64
	Method  int64
	Note    *string
}

type TransactionQR struct {
	MethodCode string  `json:"methodCode"`
	QrText     string  `json:"qrText"`
	ExpiresAt  string  `json:"expiresAt"`
	Amount     float64 `json:"amount"`
}
