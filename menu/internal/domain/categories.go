package domain

type Category struct {
	ID     int64  `json:"id,string"`
	Name   string `json:"name"`
	NameEn string `json:"nameEN"`
}
