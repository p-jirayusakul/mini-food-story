package domain

type Category struct {
	ID     int64  `json:"id,string" example:"1921144250070732800"`
	Name   string `json:"name" example:"ขนม"`
	NameEn string `json:"nameEN" example:"Dessert"`
}
