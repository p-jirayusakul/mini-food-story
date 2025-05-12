package domain

type Product struct {
	ID          int64   `json:"id,string"`
	Name        string  `json:"name"`
	NameEN      string  `json:"nameEN"`
	CategoryID  int64   `json:"categoryID,string"`
	Price       float64 `json:"price"`
	Description *string `json:"description"`
	IsAvailable bool    `json:"isAvailable"`
	ImageURL    *string `json:"imageURL"`
}

type SearchProduct struct {
	Name        string
	CategoryID  []int64
	IsAvailable bool
	OrderByType string
	OrderBy     string
	PageSize    int64
	PageNumber  int64
}

type SearchProductResult struct {
	TotalItems int64      `json:"totalItems"`
	TotalPages int64      `json:"totalPages"`
	Data       []*Product `json:"data"`
}
