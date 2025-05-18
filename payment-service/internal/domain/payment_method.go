package domain

type PaymentMethod struct {
	ID     int64  `json:"id,string"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	NameEN string `json:"name_en"`
}
