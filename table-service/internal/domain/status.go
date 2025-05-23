package domain

type Status struct {
	ID     int64  `json:"id,string" example:"1921144250070732800"`
	Code   string `json:"code" example:"ORDERED"`
	Name   string `json:"name" example:"สั่งอาหารแล้ว"`
	NameEn string `json:"nameEN" example:"Ordered"`
}
