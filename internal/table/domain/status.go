package domain

type Status struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	NameEn string `json:"nameEN"`
}
