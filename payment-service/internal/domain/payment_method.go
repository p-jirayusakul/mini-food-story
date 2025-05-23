package domain

type PaymentMethod struct {
	ID     int64  `json:"id,string" example:"1921144250070732800"`
	Code   string `json:"code" example:"PROMPTPAY"`
	Name   string `json:"name" example:"พร้อมเพย์"`
	NameEN string `json:"nameEN" example:"PromptPay"`
}
