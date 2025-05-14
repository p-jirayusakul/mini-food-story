package domain

type Status struct {
	ID     int64  `json:"id,string"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	NameEn string `json:"nameEN"`
}
