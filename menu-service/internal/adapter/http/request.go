package http

type Product struct {
	Name        string  `json:"name" validate:"required,no_special_char,max=255"`
	NameEN      string  `json:"nameEN" validate:"required,no_special_char,max=255"`
	CategoryID  string  `json:"categoryID" validate:"required,gte=1"`
	Price       float64 `json:"price"`
	Description *string `json:"description"`
	IsAvailable bool    `json:"isAvailable"`
	ImageURL    *string `json:"imageURL"`
}

type SearchMenu struct {
	PageNumber  int64   `query:"pageNumber"`
	PageSize    int64   `query:"pageSize"`
	Search      string  `query:"search" validate:"omitempty,no_special_char,max=255"`
	CategoryID  []int64 `query:"categoryID"`
	IsAvailable bool    `query:"isAvailable"`
	OrderBy     string  `query:"orderBy" validate:"omitempty,oneof=id tableNumber seats status"`
	OrderByType string  `query:"orderType" validate:"omitempty,oneof=asc desc"`
}
